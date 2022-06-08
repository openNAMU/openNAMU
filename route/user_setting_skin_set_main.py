from .tool.func import *

def user_setting_skin_set_main():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if ban_check() == 1:
            return re_error('/ban')
            
        if flask.request.method == 'POST':

            conn.commit()

            return redirect('/change')
        else:

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('user_setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">

                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))