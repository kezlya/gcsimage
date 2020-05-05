package gcsimage

//func SaveFile(fileName, contentType string, data []byte) (string, error) {
//	ctx := context.Background()
//	storageClient, err := storage.NewClient(ctx, nil)
//	if err != nil {
//		return "", err
//	}
//
//	handle, bucketErr := storageClient.Bucket(handle)
//	if bucketErr != nil {
//		return "", bucketErr
//	}
//
//	wc := handle.Object(fileName).NewWriter(ctx)
//	wc.ContentType = contentType
//
//	if _, writeErr := wc.Write(data); writeErr != nil {
//		return "", writeErr
//	}
//	if closeErr := wc.Close(); closeErr != nil {
//		return "", closeErr
//	}
//
//	return fmt.Sprintf(mediaUrl, fileName), nil
//}
