from .tool.func import *

def bbs_w_hide(bbs_num = '', post_num = ''):
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
        
        if flask.request.method == 'POST':
            pass
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('bbs_post_hide'), wiki_set(), wiki_custom(), wiki_css(['(' + bbs_name + ')' + ' (' + post_num_str + ')', 0])],
                data = render_simple_set('''
                    <form method="post">
                        <button type="submit">''' + load_lang('hide') + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str + '/' + post_num_str, load_lang('return')]]
            ))