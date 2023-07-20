# Go Cloud Storage API 

The `cloudbucket` package provides functionality to interact with a cloud storage bucket using the Google Cloud Storage service.

## Features

- File Upload: Upload a file to the cloud storage bucket.
- File Retrieval: Retrieve a file from the cloud storage bucket.
- Object Listing: List all objects in the cloud storage bucket.
- Signed URL Generation: Generate a signed URL for accessing a file in the cloud storage bucket.
- Bulk Object Retrieval: Retrieve and store all objects from the cloud storage bucket.
- Object Deletion: Delete an object from the cloud storage bucket.
- Update Object: Update an existing object in the cloud storage bucket.


## Prerequisites

Before using the `cloudbucket` package, make sure you have completed the following steps:

1. Google Cloud Storage Account: Sign up for a Google Cloud Storage account and create a bucket.

2. Service Account Key:

   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - On the burger menu on top left select **APIs & Services -> Credentials**. That opens the   page in the background. Click **CREATE CREDENTIALS** and select **Service Account**.
   - After clicking **CREATE**, you should be redirected to Create Role section. Search `Cloud Storage` and `Storage Admin`. Click on **CONTINUE**. You should be redirected to another page. Click **DONE**.
   - It creates you **Service Account**. Click on the account and you’ll be redirected to next page.
   - On Keys section. Click **Add Key” -> Create new key**.
   - In the **Create key** step, select **JSON** as the key type and click on **Create**. This will download the service account key file (`keys.json`).
   - Store the `keys.json` file securely. It will be used for authentication when accessing the cloud storage bucket.

3. Set up Authentication:

   - Ensure that you have the `keys.json` file generated in the previous step.
   - Place the `keys.json` file in a secure location within your project directory.
   - Provide the path to the `keys.json` file when initializing the `storageClient` in your code:

By following these steps, you will have a Google Cloud Storage account, a service account key (`keys.json`), and the necessary authentication set up for using the `cloudbucket` package.

## Endpoints

The `cloudbucket` package provides the following endpoints:

```
// POST METHOD
http://localhost:8080/upload-to-bucket
```
This endpoint handles file upload to the cloud storage bucket. It expects a file to be uploaded with the key "file".

```
// GET METHOD
http://localhost:8080/get-file-from-bucket/:filename
```
This endpoint retrieves a file from the cloud storage bucket. The filename should be provided as a URL parameter.

```
// GET METHOD
http://localhost:8080/list-objects-in-bucket
```
This endpoint lists all objects in the cloud storage bucket.

```
// GET METHOD
http://localhost:8080/get-signed-url/:filename
```
This endpoint generates a signed URL for accessing a specific file in the cloud storage bucket. The filename should be provided as a URL parameter.

```
// GET METHOD
http://localhost:8080/get-signed-url/:filename
```
This endpoint retrieves and stores all objects from the cloud storage bucket.

```
// DELETE METHOD
http://localhost:8080/delete-object/:filename
```
This endpoint deletes an object from the cloud storage bucket. The filename should be provided as a URL parameter.

```
// PUT METHOD
http://localhost:8080/update-object/:filename
```
This endpoints updates an existing object from the cloud storage bucket. The filename should be provided as a URL parameter

Please refer to the package documentation and code for detailed usage instructions and implementation details.

## HOW TO RUN

Clone the Repository:
```
https://github.com/Sahil-4555/go-cloud-storage-api.git
```

Navigate to the project folder:
```
cd cloud-storage-bucket-api
```

Install the required dependencies:
```
go mod tidy
```

Build and Run the application:
```
go run main.go
```

The API server will be available at:
```
http://localhost:8080
```