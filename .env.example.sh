#!/bin/sh

 # 这是docker-compose设置的“ZSW_ADMIN_PRIVATE_KEY” ENV VAR
export ZSW_CONTENT_REVIEW_PRIVATE_KEY=""

# 下面三个需要生成密钥（生成工具：https://tools.banquan.sh.cn/zsw-key-generator.html）
export KEXIN_JIEDIAN_A_PRIVATE_KEY="" 
export USER_A_PRIVATE_KEY=""
export USER_B_PRIVATE_KEY=""

# 改成你的API URL
export ZSW_API_URL="http://localhost:3031"