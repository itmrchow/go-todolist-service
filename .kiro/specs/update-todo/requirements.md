# Requirements Document

## Project Overview
1. 完成todo的更新功能 , client 可以透過API 來新增todo資料

## Requirements
1. 使用者透過client呼叫API來更新一筆todo資料
2. 使用者可更新
   1. title (必填, 20中文字以內)
   2. description (可選 , 100中文字以內)
   3. status (pending / doing / done)
   4. due date (可選 , UTC 時間 , 格式為 2006-01-02T15:04:05Z07:00)