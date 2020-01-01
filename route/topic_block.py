from .tool.func import *

def topic_block_2(conn, topic_num, num):
    curs = conn.cursor()

    topic_change_data = topic_change(topic_num)
    name = topic_change_data[0]
    sub = topic_change_data[1]

    if admin_check(3, 'blind (' + name + ' - ' + sub + '#' + str(num) + ')') != 1:
        return re_error('/error/3')

    curs.execute(db_change("select block from topic where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
    block = curs.fetchall()
    if block:
        if block[0][0] == 'O':
            curs.execute(db_change("update topic set block = '' where title = ? and sub = ? and id = ?"), [name, sub, str(num)])
        else:
            curs.execute(db_change("update topic set block = 'O' where title = ? and sub = ? and id = ?"), [name, sub, str(num)])

        rd_plus(name, sub, get_time())

        conn.commit()

    return redirect('/thread/' + str(topic_num) + '#' + str(num))