from .tool.func import *

from .edit import edit_editor

async def recent_history_add(name = 'Test', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'history_add (' + name + ')')

            today = get_time()
            content = flask.request.form.get('content', '')
            leng = '+' + str(len(content))

            history_plus(conn, 
                name,
                content,
                today,
                'Add:' + flask.request.form.get('get_ip', ''),
                flask.request.form.get('send', ''),
                leng,
                mode = 'add'
            )

            return redirect(conn, '/history/' + url_pas(name))
        else:            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'history_add'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'why') + '''" name="send">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + get_lang(conn, 'name') + '''" name="get_ip">
                        <hr class="main_hr">

                        ''' + await edit_editor(conn, ip) + '''
                    </form>
                ''',
                menu = [['history/' + url_pas(name), get_lang(conn, 'return')]]
            ))