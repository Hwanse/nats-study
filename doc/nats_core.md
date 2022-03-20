# NATS CORE 

## Message
메시지는 다음과 같은 요소들로 구성된다

- Subject 명
- Byte 배열 형태의 Payload
- 임의의 수의 헤더 필드들
- Optional로 `Reply` 주소 필드

### 추가 내용
- 메세지에는 최대 크기가 있다(nats server configuration에서 `max_payload`를 통해 설정)
- 메세지 최대 크기는 default로 1MB 크기이지만, 필요한 경우 최대 64MB까지 늘릴 수 있다
  (메세지 최대 크기는 8MB와 같이 합리적인 값을 유지하는 것이 좋음)
  
<br>

## Publish-Subscribe
- Publisher: 목적지인 `Subject`로 메시지를 보냄
- Subscriber: `Subject`를 Listening하고 있는 활성 Subscriber는 메세지를 수

### 특징
- 일대다(1:N)통신을 위한 발행-구독 메시지 배포 모델을 구현함 (일대다 모델을 Fan-out이라고도 표현함)
- Subscriber는 와일드카드(`*`, `>`)를 활용하여 다양한 형식으로 Subject를 Listening 할 수 있다


==== 그림 =====

<br>

## Request-Reply
Nats의 핵심 통신 메커니즘인 `Publish-Subscribe`를 활용해 Request-Reply 패턴을 지원한다. Requester는 특정 Subject에 `Reply Subject`와 함께 메세지를 Publish하고, Responder는 특정 Subject를 Listening하고 있다가 메세지를 수신하고 `Reply Subject`로 응답 메세지를 보낸다. 여기서 활용된 Reply Subject를 `Inbox`라고 부른다. 이는 응답자들의 위치와는 관계 없이 동적으로 Requester에게 다시 전달되는 고유한 Subject 값이다.

### 특징
- 여러 응답자가 동적 `Queue Group`을 형성할 수 있다. 또한 메세지 배포를 시작하거나 중지하기 위해서 Queue Group에서 구독자를 수동으로 추가하거나 제거할 필요가 없다. 이 과정은 자동으로 이루어지며 응답자 요소들은 필요에 따라 확장/축소가 가능하다
- Nats는 Publish-Subscribe를 기반으로 하기 때문에 요청(Request) 및 응답(Reply) 대기 시간을 측정하거나 이상 징후를 관찰한다
- Nats에서는 다중 응답도 허용하는데, 첫 번째 응답은 활용되고 시스템에서 추가적인 응답은 효율적으로 버린다. 이를 통해 정교한 패턴이 여러 응답자들을 가질 수 있고, 응답 대기시간과 Jitter를 줄일 수 있다


==== 그림 =====

<br>

## Queue Group
Nats에서는 분산 대기열(Queue)이라는 로드 밸런싱 기능을 제공한다.    
Queue Subscriber를 사용하면 시스템 내결함성을 제공하고, 워크로드 처리를 확장하는데 사용할 수 있는 Queue Subscriber 그룹 간에 메세지 전달의 균형을 맞출 수 있다 Queue 구독을 생성하기 위해 Subscriber는 Queue 이름을 등록한다.   
여기서 동일한 이름의 Queue를 가진 모든 Subscriber들이 `Queue Group`을 형성하게 된다. 이는 특별히 서버 설정이 필요하지 않으며, 지정한 Subject로 메세지를 Publish하면 `Queue Group`구성원 중 무작위로 하나가 선택되고 해당 Subscriber가 메세지를 수신하게 된다. 즉 `Queue Group`에 여러 Subscriber들이 있지만 메세지는 이중 하나에게 한번만 사용된다.

### 특징
- Queue Group 네이밍은 Subject와 동일한 규칙을 따르며, 대소문자를 구분하거나 공백 문자를 허용하지 않는다
- Queue Group의 장점은 서버의 설정을 통해서가 아닌 어플리케이션, Queue Subscriber들을 통해서 정의된다는 점이다
- Queue Group을 활용하여 서비스 확장에 이상적이다. 유연성과 구성 설정사항의 부재는 Nats가 다양한 플랫폼 기술들과 함께 작동할 수 있는 우수한 서비스 통신 기술로 만들어준다 

==== 그림 =====