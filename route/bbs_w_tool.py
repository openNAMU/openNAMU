from .tool.func import *

async def bbs_w_tool(bbs_num = '', post_num = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = ''
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        
        data += '''
            <h2>''' + get_lang(conn, 'tool') + '''</h2>
            <ul>
                <li><a href="/bbs/raw/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + get_lang(conn, 'raw') + '''</a></li>
            </ul>
        '''

        if await acl_check('', 'bbs_auth', '', '') != 1:
            curs.execute(db_change('select set_data from bbs_data where set_code = ? and set_id = ? and set_name = "pinned"'), [post_num_str, bbs_num_str])
            pinned = get_lang(conn, 'pinned') if not curs.fetchall() else get_lang(conn, 'pinned_release')

            data += '''
                <h3>''' + get_lang(conn, 'admin') + '''</h3>
                <ul>
                    <!-- <li><a href="/bbs/blind/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + get_lang(conn, 'hide') + '''</a></li> -->
                    <li><a href="/bbs/pinned/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + pinned + '''</a></li>
                </ul>
            '''

            data += '''
                <h3>''' + get_lang(conn, 'owner') + '''</h2>
                <ul>
                    <li><a href="/bbs/delete/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '">' + get_lang(conn, 'delete') + '''</a></li>
                </ul>
            '''

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'bbs_post_tool'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = data,
            menu = [['bbs/w/' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str), get_lang(conn, 'return')]]
        ))