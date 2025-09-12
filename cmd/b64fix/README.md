# Base64 Image Fix Script

This script fixes the issue where users drag and drop images directly into text fields, causing them to be stored as base64 in the database, making it heavy and slow.

## What it does

1. Scans all projects in the database
2. Finds base64 encoded images in text fields (description, problem_statement, solution, goals, cost_breakdown, impact_assessment, voluntary_contribution, feasibility)
3. Extracts the base64 images
4. Uploads them to Google Cloud Storage CDN
5. Replaces the base64 `<img>` tags with CDN URLs in the exact same position
6. Updates the database with the new content

## Usage

### First time setup
```bash
cd cmd/b64fix
go get
```

### Run the script
```bash
# From the b64fix directory
go run main.go

# Or build and run
go build -o b64fix
./b64fix
```

### Run with custom config
```bash
go run main.go
```

## Configuration

The script uses the main `config.yml` file from the project root. Make sure it contains:
- Database connection details
- GCS bucket configuration
- CDN URL
- GCS credentials file path

## Output

The script will log:
- Number of projects found
- Each base64 image found and processed
- Upload status for each image
- Final statistics (processed, uploaded, errors)

## Safety Features

- Only processes non-deleted projects
- Preserves the exact position of images in the content
- Logs all operations for audit trail
- Continues processing even if individual images fail
- Updates only changed fields

## Database Fields Processed

- `description`
- `problem_statement`
- `solution`
- `goals`
- `cost_beakdown` (note: typo in DB column name)
- `impact_assessment`
- `voluntery_contribution` (note: typo in DB column name)
- `feasibility`

## Image Storage Structure

Images are stored in GCS with the following path structure:
```
projects/{project_id}/{field_name}_{timestamp}_{index}.{extension}
```

Example:
```
projects/123e4567-e89b-12d3-a456-426614174000/description_1699123456_0.png
```