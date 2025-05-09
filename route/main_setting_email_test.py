from .tool.func import *

from .go_api_func_email import api_func_email

async def main_setting_email_test():
    with get_db_connect() as conn:
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 0)
        
        if flask.request.method == 'POST':
            render_data = await api_func_email()
            if render_data["response"] == "ok":
                data = get_lang(conn, "ok")
            else:
                data = get_lang(conn, "error")

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, "email_test"), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = data,
                menu = [["setting/external", get_lang(conn, 'return')]]
            ))
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, "email_test"), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <input name="title" placeholder="''' + get_lang(conn, "title") + '''">
                        <hr class="main_hr">
                        <input name="email" placeholder="''' + get_lang(conn, "email") + '''">
                        <hr class="main_hr">
                        <textarea  name="data" class="opennamu_textarea_500" placeholder="''' + get_lang(conn, "content") + '''"></textarea>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, "send") + '''</button>
                    </form>
                ''',
                menu = [["setting/external", get_lang(conn, 'return')]]
            ))