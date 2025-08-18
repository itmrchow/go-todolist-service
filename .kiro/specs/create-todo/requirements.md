# Requirements Document

## Project Overview
1. 完成todo的新增功能 , client 可以透過API 來新增一筆todo資料

## Requirements
1. 使用者透過client呼叫API來新增一筆todo資料
2. 使用者輸入
   1. title (必填, 20中文字以內)
   2. description (可選 , 100中文字以內)
   3. status (預設為pending , pending / doing / done)
   4. due date (可選 , UTC 時間 , 格式為 2006-01-02T15:04:05Z07:00)