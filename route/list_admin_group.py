from .tool.func import *

async def list_admin_group():
    with get_db_connect() as conn:
        curs = conn.cursor()

        list_data = '<ul>'
        org_acl_list = get_default_admin_group()

        curs.execute(db_change("select distinct name from alist order by name asc"))
        for data in curs.fetchall():
            if await acl_check('', 'owner_auth', '', '') != 1 and not data[0] in org_acl_list:
                delete_admin_group = ' <a href="/auth/list/delete/' + url_pas(data[0]) + '">(' + await get_lang("delete") + ')</a>'
            else:
                delete_admin_group = ''

            list_data += '' + \
                '<li>' + \
                    '<a href="/auth/list/add/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a>' + \
                    delete_admin_group + \
                '</li>' + \
            ''

        list_data += '' + \
            '</ul>' + \
            '<hr class="main_hr">' + \
            '<a href="/manager/8">(' + await get_lang('add') + ')</a>' + \
        ''

        return await render_template(
            await get_lang('admin_group_list'),
            list_data,
            0,
            [['manager', await get_lang('return')]]
        )
