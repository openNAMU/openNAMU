from .tool.func import *

from .api_bbs_w_post import api_bbs_w_post

def bbs_w_pinned(bbs_num = '', post_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect('/bbs/main')

        bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        if admin_check() != 1:
            return redirect('/bbs/w/' + bbs_num_str)
        
        temp_dict = json.loads(api_bbs_w_post(bbs_num_str + '-' + post_num_str).data)
        if not 'user_id' in temp_dict:
            return redirect('/bbs/main')
        
        if flask.request.method == 'POST':
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            if not curs.fetchall():
                curs.execute(db_change("insert into bbs_data (set_name, set_code, set_id, set_data) values ('pinned', ?, ?, ?)"), [post_num_str, bbs_num_str, get_time()])
            else:
                curs.execute(db_change('delete from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            
            return redirect('/bbs/w/' + bbs_num_str)
        else:
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            pinned = load_lang('pinned') if not curs.fetchall() else load_lang('pinned_release')

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('bbs_post_pinned'), wiki_set(), wiki_custom(), wiki_css(['(' + bbs_name + ')' + ' (' + post_num_str + ')', 0])],
                data = render_simple_set('''
                    <form method="post">
                        <button type="submit">''' + pinned + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str + '/' + post_num_str, load_lang('return')]]
            ))