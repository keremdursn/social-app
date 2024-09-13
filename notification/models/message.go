package models

import "gorm.io/gorm"

type MessageLog struct {
    gorm.Model
    QueueName string `json:"queue_name"`
    Message   string `json:"message"`
}