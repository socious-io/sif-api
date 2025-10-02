package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"regexp"
	"sif/src/apps/utils"
	"sif/src/config"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	database "github.com/socious-io/pkg_database"
)

type ProjectTextFields struct {
	ID                    uuid.UUID `db:"id"`
	Description           *string   `db:"description"`
	ProblemStatement      *string   `db:"problem_statement"`
	Solution              *string   `db:"solution"`
	Goals                 *string   `db:"goals"`
	CostBreakdown         *string   `db:"cost_beakdown"`
	ImpactAssessment      *string   `db:"impact_assessment"`
	VoluntaryContribution *string   `db:"voluntery_contribution"`
	Feasibility           *string   `db:"feasibility"`
}

var (
	configPath       = flag.String("c", "config.yml", "Path to the configuration file")
	base64ImageRegex = regexp.MustCompile(`<img[^>]*src="data:image/(png|jpg|jpeg|gif|webp);base64,([^"]+)"[^>]*>`)
	processedCount   = 0
	uploadedCount    = 0
	errorCount       = 0
)

func main() {
	log.Println("Starting base64 image fix process...")

	// Initialize config
	conf, err := config.Init(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	config.Config = conf

	// Initialize database
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 50,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	// Initialize uploader
	uploader := &utils.GCSUploader{
		CDNUrl:          config.Config.Upload.CDN,
		BucketName:      config.Config.Upload.Bucket,
		CredentialsFile: config.Config.Upload.Credentials,
	}

	// Process projects
	if err := processProjects(uploader); err != nil {
		log.Fatalf("Failed to process projects: %v", err)
	}

	log.Printf("Process completed! Processed: %d, Uploaded: %d, Errors: %d\n",
		processedCount, uploadedCount, errorCount)
}

func processProjects(uploader *utils.GCSUploader) error {
	ctx := context.Background()
	db := database.GetDB()

	// Fetch all projects with text fields
	query := `
		SELECT id, description, problem_statement, solution, goals, 
		       cost_beakdown, impact_assessment, voluntery_contribution, feasibility
		FROM projects 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var projects []ProjectTextFields
	if err := db.Select(&projects, query); err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	log.Printf("Found %d projects to process\n", len(projects))

	for _, project := range projects {
		if err := processProject(ctx, db, &project, uploader); err != nil {
			log.Printf("Error processing project %s: %v\n", project.ID, err)
			errorCount++
			continue
		}
		processedCount++
	}

	return nil
}

func processProject(ctx context.Context, db *sqlx.DB, project *ProjectTextFields, uploader *utils.GCSUploader) error {
	updated := false

	// Process each text field
	fields := []struct {
		name  string
		value **string
	}{
		{"description", &project.Description},
		{"problem_statement", &project.ProblemStatement},
		{"solution", &project.Solution},
		{"goals", &project.Goals},
		{"cost_breakdown", &project.CostBreakdown},
		{"impact_assessment", &project.ImpactAssessment},
		{"voluntary_contribution", &project.VoluntaryContribution},
		{"feasibility", &project.Feasibility},
	}

	for _, field := range fields {
		if *field.value == nil || **field.value == "" {
			continue
		}

		newContent, changed, err := processTextField(ctx, **field.value, project.ID, field.name, uploader)
		if err != nil {
			log.Printf("Error processing %s for project %s: %v\n", field.name, project.ID, err)
			continue
		}

		if changed {
			**field.value = newContent
			updated = true
			log.Printf("Updated %s for project %s\n", field.name, project.ID)
		}
	}

	// Update database if any field was changed
	if updated {
		if err := updateProject(db, project); err != nil {
			return fmt.Errorf("failed to update project: %w", err)
		}
		log.Printf("Successfully updated project %s in database\n", project.ID)
	}

	return nil
}

func processTextField(ctx context.Context, content string, projectID uuid.UUID, fieldName string, uploader *utils.GCSUploader) (string, bool, error) {
	matches := base64ImageRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return content, false, nil
	}

	log.Printf("Found %d base64 images in %s for project %s\n", len(matches), fieldName, projectID)

	newContent := content
	for i, match := range matches {
		fullMatch := match[0]
		imageType := match[1]
		base64Data := match[2]

		// Decode base64
		imageData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			log.Printf("Failed to decode base64: %v\n", err)
			continue
		}

		// Generate unique filename
		timestamp := time.Now().Unix()
		fileName := fmt.Sprintf("projects/%s/%s_%d_%d.%s", projectID, fieldName, timestamp, i, imageType)

		// Upload to CDN
		reader := bytes.NewReader(imageData)
		contentType := fmt.Sprintf("image/%s", imageType)

		cdnURL, err := uploader.UploadFile(ctx, fileName, contentType, reader)
		if err != nil {
			log.Printf("Failed to upload image: %v\n", err)
			continue
		}

		// Replace base64 with CDN URL
		newImgTag := fmt.Sprintf(`<img src="%s">`, cdnURL)
		newContent = strings.Replace(newContent, fullMatch, newImgTag, 1)
		uploadedCount++

		log.Printf("Uploaded image %d/%d to: %s\n", i+1, len(matches), cdnURL)
	}

	return newContent, true, nil
}

func updateProject(db *sqlx.DB, project *ProjectTextFields) error {
	query := `
		UPDATE projects 
		SET description = $2,
		    problem_statement = $3,
		    solution = $4,
		    goals = $5,
		    cost_beakdown = $6,
		    impact_assessment = $7,
		    voluntery_contribution = $8,
		    feasibility = $9,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := db.Exec(query,
		project.ID,
		project.Description,
		project.ProblemStatement,
		project.Solution,
		project.Goals,
		project.CostBreakdown,
		project.ImpactAssessment,
		project.VoluntaryContribution,
		project.Feasibility,
	)

	return err
}
