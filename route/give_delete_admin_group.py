from .tool.func import *

async def give_delete_admin_group(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if name in get_default_admin_group():
            return redirect(conn, '/auth/list')

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            curs.execute(db_change("select name from user_set where name = 'acl' and data = ? limit 1"), [name])
            if not curs.fetchall():
                await acl_check(tool = 'owner_auth', memo = 'auth list delete (' + name + ')')

                curs.execute(db_change("delete from alist where name = ?"), [name])

                return redirect(conn, '/auth/list')
            else:
                return await re_error(conn, 47)
        else:
            return easy_minify(flask.render_template(await skin_check(),
                imp = [await get_lang("delete_admin_group"), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + name + ')', 0])],
                data = '' + \
                    '<form method="post">' + \
                        '<button type="submit">' + await get_lang('delete') + '</button>' + \
                    '</form>' + \
                '',
                menu = [['auth/list', await get_lang('return')]]
            ))