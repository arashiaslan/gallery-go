import React, { useState, useCallback, useEffect, useRef } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  LinearProgress,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Typography,
} from '@mui/material';
import { Close, CloudUpload } from '@mui/icons-material';
import { useDropzone } from 'react-dropzone';
import { useUploadPhotos } from '../hooks/useBooks';
import { getErrorMessage } from '../lib/client';
import { toast } from 'sonner';

interface UploadPhotosDialogProps {
  bookId: string;
  isOpen: boolean;
  onClose: () => void;
}

interface FileWithProgress {
  file: File;
  previewUrl: string;
  progress: number;
  status: 'pending' | 'uploading' | 'success' | 'error';
  error?: string;
}

// Must match the backend MAX_UPLOAD_MB (25) default
const MAX_FILE_SIZE = 25 * 1024 * 1024;

export const UploadPhotosDialog: React.FC<UploadPhotosDialogProps> = ({
  bookId,
  isOpen,
  onClose,
}) => {
  const [files, setFiles] = useState<FileWithProgress[]>([]);
  const uploadMutation = useUploadPhotos();

  const onDrop = useCallback((acceptedFiles: File[]) => {
    const newFiles = acceptedFiles.map((file) => ({
      file,
      previewUrl: URL.createObjectURL(file),
      progress: 0,
      status: 'pending' as const,
    }));
    setFiles((prev) => [...prev, ...newFiles]);
  }, []);

  const filesRef = useRef<FileWithProgress[]>(files);
  useEffect(() => {
    filesRef.current = files;
  }, [files]);
  useEffect(() => {
    return () => {
      filesRef.current.forEach((f) => URL.revokeObjectURL(f.previewUrl));
    };
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    maxSize: MAX_FILE_SIZE,
    maxFiles: 20,
    accept: {
      'image/*': ['.jpeg', '.jpg', '.png', '.webp'],
    },
    onDropRejected: (rejections) => {
      const tooLarge = rejections.some((r) =>
        r.errors.some((e) => e.code === 'file-too-large'),
      );
      if (tooLarge) {
        toast.error('File too large: max 25MB per photo');
      }
    },
  });

  const releasePreviews = (list: FileWithProgress[]) => {
    list.forEach((f) => URL.revokeObjectURL(f.previewUrl));
  };

  const handleUpload = async () => {
    if (files.length === 0) {
      toast.error('Select at least 1 photo');
      return;
    }

    const filesToUpload = files.map((f) => f.file);

    try {
      setFiles((prev) => prev.map((f) => ({ ...f, status: 'uploading' as const })));

      await uploadMutation.mutateAsync({
        bookId,
        files: filesToUpload,
        onUploadProgress: (progress) => {
          setFiles((prev) => prev.map((f) => ({ ...f, progress })));
        },
      });

      toast.success(`${filesToUpload.length} photos uploaded`);
      releasePreviews(files);
      setFiles([]);
      onClose();
    } catch (error) {
      toast.error(getErrorMessage(error, 'Upload failed'));
      setFiles((prev) =>
        prev.map((f) => ({ ...f, status: 'error' as const, error: 'Upload failed' })),
      );
    }
  };

  const removeFile = (index: number) => {
    setFiles((prev) => {
      const target = prev[index];
      if (target) URL.revokeObjectURL(target.previewUrl);
      return prev.filter((_, i) => i !== index);
    });
  };

  const batchProgress =
    files.length > 0
      ? Math.round(files.reduce((sum, f) => sum + f.progress, 0) / files.length)
      : 0;

  return (
    <Dialog
      open={isOpen}
      onClose={() => {
        if (!uploadMutation.isPending) onClose();
      }}
      maxWidth="xs"
      fullWidth
    >
      <DialogTitle>Upload photos</DialogTitle>

      <DialogContent>
        {files.length === 0 ? (
          <Box
            {...getRootProps()}
            sx={{
              border: '2px dashed',
              borderColor: isDragActive ? 'primary.main' : 'divider',
              borderRadius: 2,
              py: 5,
              px: 2,
              textAlign: 'center',
              cursor: 'pointer',
              bgcolor: isDragActive ? 'action.hover' : 'transparent',
              '&:hover': { bgcolor: 'action.hover' },
            }}
          >
            <input {...getInputProps()} />
            <CloudUpload sx={{ fontSize: 36, color: 'text.secondary' }} />
            <Typography variant="body2" sx={{ mt: 1 }}>
              {isDragActive ? 'Drop the photos here...' : 'Drag photos here or click to select'}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              JPG, PNG, WebP (max 25MB)
            </Typography>
          </Box>
        ) : (
          <>
            <List dense disablePadding>
              {files.map((file, index) => (
                <ListItem
                  key={file.previewUrl}
                  secondaryAction={
                    file.status !== 'uploading' ? (
                      <IconButton
                        edge="end"
                        size="small"
                        onClick={() => removeFile(index)}
                        aria-label={`Remove ${file.file.name}`}
                      >
                        <Close fontSize="small" />
                      </IconButton>
                    ) : null
                  }
                  sx={{ px: 0 }}
                >
                  <ListItemAvatar sx={{ minWidth: 56 }}>
                    <Box
                      component="img"
                      src={file.previewUrl}
                      alt=""
                      sx={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 1 }}
                    />
                  </ListItemAvatar>
                  <ListItemText
                    primary={file.file.name}
                    secondary={
                      file.status === 'uploading'
                        ? 'Uploading...'
                        : file.status === 'error'
                          ? file.error
                          : `${(file.file.size / 1024 / 1024).toFixed(2)} MB`
                    }
                    slotProps={{
                      primary: { noWrap: true, variant: 'body2' },
                      secondary: {
                        sx: { color: file.status === 'error' ? 'error.main' : undefined },
                      },
                    }}
                  />
                </ListItem>
              ))}
            </List>

            {uploadMutation.isPending && (
              <LinearProgress variant="determinate" value={batchProgress} sx={{ mt: 1 }} />
            )}

            {!uploadMutation.isPending && (
              <Button
                fullWidth
                sx={{ mt: 1.5 }}
                onClick={() => {
                  const input = document.createElement('input');
                  input.type = 'file';
                  input.multiple = true;
                  input.accept = 'image/*';
                  input.onchange = (e) => {
                    const fileList = (e.target as HTMLInputElement).files;
                    if (fileList) {
                      onDrop(Array.from(fileList));
                    }
                  };
                  input.click();
                }}
              >
                Add more photos
              </Button>
            )}
          </>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={uploadMutation.isPending}>
          Cancel
        </Button>
        <Button
          variant="contained"
          disableElevation
          onClick={handleUpload}
          disabled={files.length === 0 || uploadMutation.isPending}
        >
          {uploadMutation.isPending ? <CircularProgress size={18} color="inherit" /> : 'Upload'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
