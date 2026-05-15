# 1. 克隆项目
#cd Goto_Web3

# 2. 配置环境变量
# export DB_DSN="username:password@tcp(127.0.0.1:3306)/goto_web3?parseTime=true&charset=utf8mb4"
# export JWT_SECRET="your-secret-key"
# export PORT=8080

# 3. 导入学习计划数据
cd backend && go run cmd/seed/main.go ../sources/web3_infra_3month_plan.md

# 4. 启动服务
go run cmd/server/main.go
# 访问 http://localhost:8080

cd ../