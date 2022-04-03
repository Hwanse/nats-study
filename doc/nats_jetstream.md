# NATS JetStream

기본 Nats Core 기능에 더 높은 서비스 품질을 가능하게 하기 위해서 JetStream이라는 내장형 분산
지속성 시스템. JetStream은 스트리밍 기술의 복잡성, 취약성, 확장성의 부족과 같은 문제들을 해결
하기 위해 만들어졌다. 

- `nats-server` 모듈에 내장 되어 있음
- Nats 서버중 1개만 있으면 JetStream 기능을 사용할 수 있다

# JetStream에 의해서 활성화된 기능

## Streaming: Publisher와 Subscriber 간에 일시적 decoupling

기본 Pub/Sub 메세징의 테넌트 중 하나는 Publisher와 Subscriber들 간 사이에 일시적 결합이 생긴다. 메시징 시스템이 Publisher와 Subscriber 사이의 시간적 decoupling을 제공하는 전통적인 방법은 
`durable subscriber`기능을 통해서 또는 `queues`를 활용하는 것이지만 어느 한쪽도 완벽한 방법은아니다.

- 메세지가 Publish 되기 전에 durable subscriber를 만들어야 한다
- queue는 로드 밸런싱 또는 소비성을 위한 기능이고 메세지 재생을 위한 메커니즘으로 사용되지 않는다

그러나 최근 이러한 decoupling을 제공하기 위해 streaming이 고안되어 주류가 되었다. `Streams`은 하나 이상의 Subject에 게시된 메세지를 캡처 및 저장하고, 클라이언트 어플리케이션이 `Stream`에 저장된 메세지 전체 또는 일부를 `reply`(or consume)하기 위해 언제든지 subscriber들을 생성할 수 있도록 한다

### Replay Policies

Jetsream consumer는 소비하는 어플리케이션 측에 여러 replay policies를 지원

- 현재 Stream에 저장된 모든 메세지, 완전한 replay를 의미하고 replay policy(재생 속도)를 설정할 때 다음과 같은 선택지가 있다
  - instant: 메세지가 최대한 빠르게 consumer에게 전달되는 정책
  - original: Stream에 게시된 속도로 메세지가 consumer에게 전달됨, 이 옵션은 운영 환경 트래픽에서 유용할 수 있음
- Stream에 저장된 마지막 메세지 또는 각 Subject에 대한 마지막 메세지
- 특정 '시퀀스 번호'에서 부터 시작
- 특정 '시간'에서 부터 시작

### Retention policies 및 limits

기본 Nats Core 기능 위에 새로운 기능과 더 높은 서비스 품질을 가능하게 한다. 사실 실직적으로는 Stream을 항상 무한 연장할 수가 없으므로 JetStream은 Stream에 크기 제한 및 여러 보존 정책을 지원. Stream에는 다음과 같이 limit을 부여할 수 있다.

**Limits**

- 메세지 최대 유효 기한
- Stream의 총 크기(byte 단위)
- Stream의 최대 메세지의 개수
- 각 메세지별 최대 크기
- 삭제 정책 지정 가능: 한도에 도달한 상태에서 Stream에 새로운 메세지가 Publish 되면 해당 새매세지를 위한 공간을 만들려고 현재 Stream의 가장 오래된 메세지 또는 가장 최신 메세지를 버리도록 설정 가능
- 특정 시점에 Stream에 대한 최대 consumer의 숫자를 제한할 수 있다

**Retention policy**

다음 유형의 Stream retention policy 선택 가능

- limits(default)
- interest: Stream에 소비자가 있는 한 메세지는 Stream에서 유지 됨
- work queue: Stream은 공유 Queue로 사용되고 메세지는 소비될 때 제거된다

### Persistent distributed storage

필요에 따라 메세지 저장소의 내구성과 복원력을 선택 할 수 있다

- Memory storage
- File stroage
- 내결함성을 위해 nats server 간에 복제

JetStream에서 메세지를 저장하기 위한 구성 방식은 메세지가 사용되는 방식과 별도로 정의된다.
저장소는 `Stream`에 정의되고 메세지 소비는 여러 `Consumer`들에 의해 정의된다

## De-Coupled flow control

JetStream은 Stream에 대한 분리된 흐름 제어를 제공한다. 흐름제어는 Publisher가 모든 Consumer 중 가장 느린 속도로 Publish 하도록 제한되지 않고, 각 클라이언트 애플리케이션(Publisher or Consumer)과 Nats 서버 간에 개별적으로 발생한다.
JetStream Publish 호출을 Stream에 Publish하는데 사용할 경우 Publisher와 Nats 서버간에 승인 메커니즘이 있으며, 동기식 또는 비동깃직 JetStream Publish를 선택할 수 있다

Subscriber 측에서는 Stream으로 부터 메세지를 수신하거나 소비하는 클라이언트 애플리케이션으로 Nats 서버로 부터의 메세지 전송도 제어된다


## Consumers

JetStream의 `Consumer`는 Stream의 뷰어이며, Stream에 저장된 메세지의 복사본을 수신하기 위해 클라이언트 애플리케이션에 의해 등록된다.

### Fast Push Consumer

클라이언트 애플리케이션은 승인되지 않는 fast push consumer를 사용하여 지정된 subject 또는 받은 Inbox에 대한 메세지를 가능한 한 빨리 수신하도록 선택할 수 있다. 이 consumer는 Stream에서 메세지를 소비보다는 'replay'를 목적으로 사용된다

### 배치를 통한 수평 확장성이 뛰어난 pull consumer

클라이언트 애플리케이션은 요구 주도형이며 일괄 처리를 지원하며 메시지 수신 및 처리를 명시적으로 승인해야 하는 풀 컨슈머를 사용하여 공유할 수도 있습니다. 즉, 스트림에서 메시지를 처리할 뿐만 아니라 소비(즉, 스트림의 분산 큐로서 사용)할 수 있습니다.
풀 컨슈머는 (예를 들어) 파티션을 정의하거나 폴트 톨러런스에 대해 걱정할 필요 없이 스트림 내의 메시지 처리 또는 소비를 쉽고 투과적으로 수평적으로 확장하기 위해 (큐 그룹과 마찬가지로) 애플리케이션 간에 공유할 수 있습니다.
주의: 사용자의 Fetch 호출에 (합리적인) 타임아웃을 전달하여 루프에서 호출할 수 있으므로 풀 컨슈머를 사용한다고 해서 업데이트(스트림에 게시된 새 메시지)를 실시간으로 응용 프로그램에 '푸시'할 수 없는 것은 아닙니다.

### consumer 확인 사항

가능한 가장 빠른 메시지 전달을 위해 서비스 품질을 거래하는 승인되지 않은 소비자를 사용하도록 결정할 수 있지만 대부분의 처리는 멱등성이 아니며 더 높은 서비스 품질(예: 일부 메시지에서 처리되지 않거나 두 번 이상 처리됨) 확인된 소비자를 사용하려고 할 것입니다. JetStream은 두 가지 이상의 승인 유형을 지원합니다.

- 일부 소비자는 확인 중인 메시지의 시퀀스 번호까지 모든 메시지를 확인하는 것을 지원하고, 일부 소비자는 최고 품질의 서비스를 제공하지만 각 메시지의 수신 및 처리를 명시적으로 확인하고 서버가 대기할 최대 시간을 확인해야 합니다. (소비자에게 연결된 다른 프로세스에) 다시 전달하기 전에 특정 메시지에 대한 승인
- 부정적인 승인을 다시 보낼 수도 있습니다
- 진행 중인 확인을 보낼 수도 있습니다(문제의 메시지를 아직 처리 중이고 확인 또는 확인하기 전에 시간이 더 필요함을 나타내기 위해)

## K/V Store
JetStream은 지속성 계층이며 스트리밍은 해당 계층 위에 구축된 기능 중 하나일 뿐입니다.

또 다른 기능(일반적으로 메시징 시스템에서 사용할 수 없거나 심지어 메시징 시스템과 연결되어 있지도 않음)은 JetStream 키/값 저장소입니다. 키와 관련된 값 메시지를 저장, 검색 및 삭제하고 해당 키에 발생하는 변경 사항을 감시(수신)하고 심지어 특정 키에서 발생한 값(및 삭제)의 기록을 검색합니다.