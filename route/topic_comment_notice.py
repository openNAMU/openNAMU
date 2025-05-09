from .tool.func import *

async def topic_comment_notice(topic_num = 1, num = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        topic_num = str(topic_num)
        num = str(num)
        
        if await acl_check(tool = 'toron_auth', memo = 'notice (code ' + topic_num + '#' + num + ')') == 1:
            return await re_error(conn, 3)

        curs.execute(db_change("select code from topic where code = ? and id = ?"), [topic_num, num])
        if curs.fetchall():
            curs.execute(db_change("select top from topic where code = ? and id = ?"), [topic_num, num])
            top_data = curs.fetchall()
            if top_data:
                if top_data[0][0] == 'O':
                    curs.execute(db_change("update topic set top = '' where code = ? and id = ?"), [topic_num, num])
                else:
                    curs.execute(db_change("update topic set top = 'O' where code = ? and id = ?"), [topic_num, num])

            do_reload_recent_thread(conn, 
                topic_num, 
                get_time()
            )

        return redirect(conn, '/thread/' + topic_num + '#' + num)
