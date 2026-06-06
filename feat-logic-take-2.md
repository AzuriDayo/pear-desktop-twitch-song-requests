# SR Logic take 2

all uses queue mutex

Twitch events:

- on !sr
  - validate request and add song to queue
    - if queue was empty, instantly run add song to next in queue pear-desktop

Pear events:

- on new song playing type VIDEO_CHANGED
  - validate currently playing song is the one on the head of queue
    - pop head, save in memory requester
    - if queue is not empty after pop, add song to next in queue pear-desktop
