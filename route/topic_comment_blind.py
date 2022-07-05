from .tool.func import *

def topic_comment_blind(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        topic_num = str(topic_num)
        num = str(num)
        
        if admin_check(3, 'blind (code ' + topic_num + '#' + num + ')') != 1:
            return re_error('/error/3')

        curs.execute(db_change("select block from topic where code = ? and id = ?"), [topic_num, num])
        block = curs.fetchall()
        if block:
            if block[0][0] == 'O':
                curs.execute(db_change("update topic set block = '' where code = ? and id = ?"), [topic_num, num])
            else:
                curs.execute(db_change("update topic set block = 'O' where code = ? and id = ?"), [topic_num, num])

            do_reload_recent_thread(
                topic_num, 
                get_time()
            )

            conn.commit()

        return redirect('/thread/' + topic_num + '#' + num)