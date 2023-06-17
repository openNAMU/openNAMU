from .tool.func import *

def bbs_make():   
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/3')
        
        if flask.request.method == 'POST':
            
            curs.execute(db_change('select set_id from bbs_set where set_name = "bbs_name" order by set_id + 0 desc'))
            db_data = curs.fetchall()

            bbs_num = str(int(db_data[0][0]) + 1) if db_data else '1'
            bbs_name = flask.request.form.get('bbs_name', 'test')
            bbs_type = flask.request.form.get('bbs_type', 'comment')
            bbs_type = bbs_type if bbs_type in ['comment', 'thread'] else 'comment'

            curs.execute(db_change("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_name', '', ?, ?)"), [bbs_num, bbs_name])
            curs.execute(db_change("insert into bbs_set (set_name, set_code, set_id, set_data) values ('bbs_type', '', ?, ?)"), [bbs_num, bbs_type])

            conn.commit()

            return redirect('/bbs/main')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('bbs_make'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + load_lang('bbs_name') + '''" name="bbs_name">
                        <hr class="main_hr">
                        
                        <select name="bbs_type">
                            <option value="comment">''' + load_lang('comment_base') + '''</option>
                            <option value="thread">''' + load_lang('thread_base') + '''</option>
                        </select>
                        <hr class="main_hr">
                        
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['bbs/main', load_lang('return')]]
            ))