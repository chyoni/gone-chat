# Gone-Chat

## Index

- #01 Init

- #02 Chat Server Connection

- #03 Read and Write Message 1

- #04 Global Chat Server DONE

- #05 Private rooms chat 1

- #06 Database Init 1

- #07 Create User 1

- #08 Create User 2 (Hashed Password)

- #09 Create Room 1

- #10 Login 1

- #11 Login 2 (Generate JWT)

  - go get github.com/dgrijalva/jwt-go

- #12 Login 3 (redis, uuid)

  - go get github.com/go-redis/redis/v7
  - go get github.com/twinj/uuid

  > 현재 구현된 JWT는 만료 시간이 지나야만 토큰이 유효하지 않게 되고 그렇게 되면
  > 유저가 로그인 후 몇 분 안지나서 로그아웃했을 때 해당 토큰이 그대로 살아있기 때문에 보안적 위험이 있다.
  > 이를 해결하기 위해, redis라는 key-value storage를 사용한다.
  > redis에서 key는 unique한 값이 필요하기 때문에 uuid를 사용하여 key를 unique하게 generate한다.
  > redis의 key는 uuid로 되어 있기 때문에 여러 디바이스에서 로그인이 가능하고 각 디바이스의 token이 별개로 생성된다.
  > 이 말은, 디바이스 별 로그아웃이 가능하다는 말이고 멋진 일이다.

  > cli로 redis server를 키는 방법은 'brew services start redis'
  > cli로 redis를 실행하는 방법은 'redis-cli'
  > 'keys \*'를 입력하면 redis에 저장된 모든 key들을 출력
  > 'get [key 이름]'를 입력하면 '[key 이름]'이라는 key의 value를 출력

- #13 Backend Temporary end

- #14 Update User Alias

- #15 Update User Password

  > should know why beforeupdate hook have limitation, then so i written hashed password code on api func

- #16 Delete User

- #17 Get Me

- #18 AllowedHeaders

  > Origin이 같지 않을 때, 허용한 Origin에서 추가하는 Headers를 허용하지 않으면 이 또한 CORS Error가 발생
  > 그렇기 때문에 허용한 Origin에서 추가하는 Headers정보 중 허용하는 Key를 설정해줘야 한다.

- #19 Token Invalid Response

- #20 Upload Image to AWS S3

  - go get -u github.com/aws/aws-sdk-go/...

- #21 Update Database avatar with uploaded image

- #22 Get Rooms

- #23 Get Users By Room

- #24 sync.Once for database
