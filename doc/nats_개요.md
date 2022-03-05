
# NATS

NATS는 오픈 소스 메시징 시스템(메시지 지향 미들웨어)이다. NATS 서버는 Go 프로그래밍 언어로 작성되었다. 
서버와의 인터페이스를 위한 클라이언트 라이브러리는 주요 프로그래밍 언어로 이용이 가능하다.     
NATS의 핵심 설계 원리는 성능, 확장성, 쉬운 이용이다

## Subject 기반의 메세징
기본적으로 NATS는 `Subject`에 의존적인 메세지 publish(게시)/listening(수신)을 다루는 내용이다. 

**Subject**
Subject는 간단하게 Publisher와 Subscriber가 서로를 찾기 위해 사용되는 이름(or 주소)을 형성하는 문자열이다.

- Subject Naming을 다음 규칙 준수를 권장
  - a~z, A~Z, 0~9 문자들의 집합(공백 없음)  
  - 허용 특수문자: `.`(subject에서 토큰 구분 목적), `*` or `>`(와일드카드 사용 목적)
  - 토큰: 토큰은 `.`를 구분자로 두어 각각 사이사이에 위치한 단어들을 지칭

- Subject는 계층형 구조 표현이 가능하다
  - 계층을 표현하기 위해 구분자로 `.`를 사용
  - 예시
    ```text
    time.us
    time.us.east
    time.us.east.atlanta
    time.eu.east
    time.eu.warsaw
    ```

**와일드카드**
NATS는 `.`로 구분된 Subject에서 하나 이상의 요소를 대신할 수 있는 두 가지의 와일드카드를 제공한다.(`*` or `>`)
Subscriber는 와일드카드를 사용하여 한 번의 구독으로 여러 Subject의 수신이 가능해지지만, Publisher는 와일드카드 없이
정확하게 명시된 Subject를 사용해야한다.

- `*`: `*`는 단일 토큰 매칭 목적으로 사용
- `>`: `>`는 복수개의 토큰을 매칭하기 위해 사용
- `*`와 `>` 혼합 사용도 가능하다
- Subject의 최대 토큰 수를 최대 16개까지만 사용할 것을 권장

<br>

# NATS Module 2가지

## Core NATS
Core NATS는 NATS 서비스에서 제공하는 기본 기능 및 서비스 품질의 집합이다. 이는 JetStream 기능이 활성화된 
`nats-server` 인스턴스가 없는 상태를 말한다  

**Core NATS 메세징 모델**

- Publish-Subscribe
- Request-Reply
- Queue Groups

## JetStream
Nats Server 2.2 버전부터 지원하는 기능사항이며, Core NATS 기능과 동시에 더 높은 수준의 메시징 보장하기 위한 기능사항이다.
JetStream `nats-server` 인스턴스에 내장되어 있으며 NATS 서버 1개만 있으면 사용이 가능하다.

**JetStream의 기능사항**

- Streams
- Consumers
- Key/Value Store(K/V Store)