from .tool.func import *

def bbs_delete(bbs_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect('/bbs/main')
        
        bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)

        if admin_check() != 1:
            return redirect('/bbs/w/' + bbs_num_str)
        
        if flask.request.method == 'POST':
            curs.execute(db_change('delete from bbs_data where set_id = ?'), [bbs_num_str])
            curs.execute(db_change('delete from bbs_set where set_id = ?'), [bbs_num_str])
            curs.execute(db_change('delete from bbs_data where set_id like ?'), [bbs_num_str + '-%'])
            
            return redirect('/bbs/w/' + bbs_num_str)
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('bbs_delete'), wiki_set(), wiki_custom(), wiki_css(['(' + bbs_name + ')', 0])],
                data = render_simple_set('''
                    <form method="post">
                        <span>''' + load_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                '''),
                menu = [['bbs/set/' + bbs_num_str, load_lang('return')]]
            ))