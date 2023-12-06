from .tool.func import *

def bbs_delete(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect('/bbs/main')
        else:
            bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)

        if admin_check() != 1:
            return redirect('/bbs/w/' + bbs_num_str)
        
        if flask.request.method == 'POST':
            pass
        else:
            pass