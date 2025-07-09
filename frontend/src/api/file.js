import http from './http'

// 上传文件
export const uploadFile = (formData, onProgress) => {
  return http.post('/files/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress) {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        )
        onProgress(percentCompleted)
      }
    }
  })
}

// 获取文件信息
export const getFileInfo = (fileId) => {
  return http.get(`/files/${fileId}`)
}

// 删除文件
export const deleteFile = (fileId) => {
  return http.delete(`/files/${fileId}`)
}

// 获取文件下载URL（带认证token）
export const getDownloadUrl = (fileId) => {
  const token = localStorage.getItem('token')
  return `/api/files/${fileId}/download?token=${token}`
}

// 获取文件预览URL（带认证token）
export const getPreviewUrl = (fileId) => {
  const token = localStorage.getItem('token')
  return `/api/files/${fileId}/preview?token=${token}`
}

// 获取静态文件URL
export const getStaticUrl = (filePath) => {
  return `/uploads/${filePath}`
} 