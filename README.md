# finpc-example

- 폴더별 기능
.github/workflows: 
.local: docker-compose, PostgreSQL 5432
client: next.js 부분
proxy: aws
server: golang gRPC
terraform: aws cloud


github에서 action(proxy, client, server 3개 다 배포) 후
$terraform output 하면 endpoint = ~~~~~ 라고 뜰거다


_raw postgres_password




*terraform
- main.tf
    - state: local file (terraform.tfstate 생성됨)
    - provider: resourse, data 제어
    - 
- variables.tf에 변수들
- outputs.tf: output 명령시 실행되는 것
- 선언
    - resource: create 
    - data: read

