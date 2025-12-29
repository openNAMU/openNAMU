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
            return await render_template(
                await get_lang("delete_admin_group"),
                '' + \
                    '<form method="post">' + \
                        '<button type="submit">' + await get_lang('delete') + '</button>' + \
                    '</form>' + \
                '',
                '(' + name + ')',
                [['auth/list', await get_lang('return')]]
            )
