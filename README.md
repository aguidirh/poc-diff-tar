# Diff Tar POC

This is a POC for the diff mechanism required to extract only images that were not mirrored previously in oc-mirror

## Using the command line
Run the command below specifying in the first argument the history and in the second a file which contains the images to be mirrored

```
go run main.go .MirroringHistory .newMirroring
```

In the root folder, there are two hidden files:

``.MirroringHistory`` - it is a hidden file which simulates the file that stores all digests mirrored to a target registry.

``.newMirroring`` - it is a hidden file which simulates the images that are being prepared for a mirroring.