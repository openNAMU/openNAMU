from .tool.func import *

from .edit_delete import edit_delete

def edit_delete_multiple():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if admin_check() != 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if do_edit_send_check(send) == 1:
                return re_error('/error/37')
            
            if do_edit_text_bottom_check_box_check(agree) == 1:
                return re_error('/error/29')
            
            all_title = re.findall(r'([^\n]+)\n', flask.request.form.get('content', '').replace('\r', '') + '\n')
            for name in all_title:
                edit_delete(name)

            return redirect('/recent_change')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('many_delete'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <textarea class="opennamu_textarea_500" placeholder="''' + load_lang('many_delete_help') + '''" name="content"></textarea>
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['manager/1', load_lang('return')]]
            ))