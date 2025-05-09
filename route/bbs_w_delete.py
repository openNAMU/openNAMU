from .tool.func import *

from .go_api_bbs_w import api_bbs_w

async def bbs_w_delete(bbs_num = '', post_num = '', comment_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/bbs/main')

        bbs_name = db_data[0][0]
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)

        if await acl_check('', 'owner_auth', '') == 1:
            return redirect(conn, '/bbs/in/' + bbs_num_str)
        
        temp_dict = await api_bbs_w(bbs_num_str + '-' + post_num_str)
        if not 'user_id' in temp_dict:
            return redirect(conn, '/bbs/main')
        
        if flask.request.method == 'POST':
            if comment_num == '':
                curs.execute(db_change('delete from bbs_data where set_code = ? and set_id = ?'), [post_num_str, bbs_num_str])
                curs.execute(db_change('delete from bbs_set where set_code = ? and set_id = ?'), [post_num_str, bbs_num_str])
                curs.execute(db_change('delete from bbs_data where set_id = ? or set_id like ?'), [bbs_num_str + '-' + post_num_str, bbs_num_str + '-' + post_num_str + '-%'])
                
                return redirect(conn, '/bbs/in/' + bbs_num_str)
            else:
                comment_num_split = comment_num.split('-')
                
                set_id = bbs_num_str + '-' + post_num_str
                set_id_sub = '-'.join(comment_num_split[:-1])
                if set_id_sub != '':
                    set_id += '-' + set_id_sub

                set_code = comment_num_split[len(comment_num_split) - 1]

                curs.execute(db_change("update bbs_data set set_data = '' where set_code = ? and set_id = ?"), [set_code, set_id])
                
                return redirect(conn, '/bbs/w/' + bbs_num_str + '/' + post_num_str)
        else:
            sub = '(' + bbs_name + ')'
            sub += ' (' + post_num_str + ')'
            
            name = get_lang(conn, 'bbs_comment_delete')
            if comment_num == '':
                name = get_lang(conn, 'bbs_post_delete')
            else:
                sub += ' (' + comment_num + ')'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
                data = render_simple_set(conn, '''
                    <form method="post">
                        <span>''' + get_lang(conn, 'delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                '''),
                menu = [['bbs/w/' + bbs_num_str + '/' + post_num_str, get_lang(conn, 'return')]]
            ))