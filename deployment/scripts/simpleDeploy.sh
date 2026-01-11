cd ~
rm -rf full-stack-cloud-test
git clone https://github.com/Henri-lab/full-stack-cloud-test.git
cd ~/full-stack-cloud-test

echo "正在准备部署..."
read -r -p "准备就绪，按 Enter 开始部署: "
echo "开始部署..."

cp .env.production deployment/.env
docker-compose -f deployment/docker-compose.prod.yml up -d --build
