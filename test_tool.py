# Load
import time
import os
import platform
import urllib
import zipfile

from route.tool.func import *

data_db_set = class_check_json()

db_data_get(data_db_set['type'])
do_db_set(data_db_set)

load_db = get_db_connect()

conn = load_db.__enter__()
curs = conn.cursor()

print('1. Add virtual doc')

what_i_do = input('Select : ')
if what_i_do == '1':
    doc_count = int(input('Count : '))
    
    for for_a in range(doc_count):
        name = 'test_' + str(for_a)
        content = '''[include(틀:주요 문서)]
[include(틀:위키 엔진)]

||||<tablealign=right><tablebordercolor=#008679><#008679> {{{#white {{{+2 오픈나무
''\'openNAMU\'''}}}}}} ||
|||| [[파일:오픈나무 로고.png|width=200px]] ||
||<colbgcolor=#008679><colcolor=white> 현재 상태 || ~~계속~~ 개발 중 ||
|||| 최신 빌드 ||
|| stable || v3.4.5 (stable1) (beta3) (dev14) ||
|| beta || v3.4.6-RC1 (stable1) (beta8-107) ||
|| dev || v3.4.6-RC1 (stable1) (beta8-107)[* 거의 매일 바뀜] ||
|||| 역사 ||
|| nodeJS || 2016-04-23 ||
|| Python || 2017-01-06 ||
|| 링크 || [[https://github.com/openNAMU/openNAMU|깃허브]]
[[https://github.com/openNAMU|스킨, 서브 자료들]] ||
[목차]
[clearfix]
== 개요 ==
[[나무마크]]를 ~~대충~~ 지원하는 [[파이썬]] 위키 엔진임다.

== 왜 오픈나무임? ==
[[basix|원래 개발자]]가 [[더 시드]] 같은 엔진을 만들겠다라는 명목하에 open[[나무위키|NAMU]]라고 지었습니다. 

== 왜 만듬? ==
과거에는 저도 [[대한위키실록|위키]]를 열어봤던 사람으로써 [[도쿠위키]]와 [[모니위키]], [[미디어위키]]를 쓰다 불편한 점을 개선하기 위해서 새로 만들었슴다.

[[미디어위키]]는 정말 좋은 엔진이지만 아무래도 한국적 특성에 뭔가 안 맞는 것 같아서 이걸 만들었슴다.

참고로 쓰다보면 이런 게 왜 있지? 싶은 요소가 있는데[* 예를 들면 틀 링크] 그건 그냥 제가 쓰려고 만들었슴다.

=== 영향을 많이 받은 위키 관련 요소 ===
[include(틀:오픈나무 개발 이유)]

== 역사 ==
 * [[/역사]]
  * [[/v3.4.6]]
== 개발 이념 ==
 * [[/개발 이념]]
== 주의 ==
 * 프로덕션 용으로는 사용 안하는 걸 권장 드립니다. ~~[[시한폭탄|언제 터질 지 모릅니다]]~~
== 사용하는 위키 ==
 * [[/사용 위키]]

== 나머지 개발자의 헛소리 적는 공간 ==
 * [[/연구]]
 * [[/차기 계획]]
 * [[/개발 현황]]
 * [[/반성]]

[[분류:오픈나무]]'''
        today = get_time()
        send = 'test'
        ip = '127.0.0.1'
        leng = '0'

        curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, content])
        print(for_a)

        test_case = [[
            for_b,
            name,
            content,
            today,
            ip,
            send,
            leng,
            ''
        ] for for_b in range(1, 151)]
        curs.executemany(db_change(
            "insert into history (id, title, data, date, ip, send, leng, hide, type) " + \
            "values (?, ?, ?, ?, ?, ?, ?, '', ?)"
        ), test_case)

    conn.commit()