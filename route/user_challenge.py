from .tool.func import *

def do_make_challenge_design(img, title, info, disable = 0):
    if disable == 1:
        table_style = 'style="border: 2px solid green"'
    else:
        table_style = 'style="border: 2px solid red"'
    
    return '''
        <table id="main_table_set" ''' + table_style + '''>
            <tr>
                <td id="main_table_width_quarter" rowspan="2">
                    <span style="font-size: 64px;">''' + img + '''</span>
                </td>
                <td>
                    <span style="font-size: 32px;">''' + title + '''</span>
                </td>
            </tr>
            <tr>
                <td>''' + info + '''</td>
        </table>
        <hr class="main_hr">
    '''

def user_challenge():    
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        ip = ip_check()
        if ip_or_user(ip) == 1:
            return redirect('/user')

        if flask.request.method == 'POST':
            pass
        else:
            data_html_green = ''
            data_html_red = ''
            
            data_html_green += do_make_challenge_design(
                '🌳',
                load_lang('challenge_title_register'), 
                load_lang('challenge_info_register', 1),
                1
            )
            
            curs.execute(db_change('select count(*) from history where ip = ?'), [ip])
            db_data = curs.fetchall()
            
            disable = 1 if db_data[0][0] >= 1 else 0
            data_html = do_make_challenge_design(
                '🔰',
                load_lang('challenge_title_first_contribute'), 
                load_lang('challenge_info_first_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 10 else 0
            data_html = do_make_challenge_design(
                '📝',
                load_lang('challenge_title_tenth_contribute'), 
                load_lang('challenge_info_tenth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 100 else 0
            data_html = do_make_challenge_design(
                '🖊️',
                load_lang('challenge_title_hundredth_contribute'), 
                load_lang('challenge_info_hundredth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 1000 else 0
            data_html = do_make_challenge_design(
                '🏅',
                load_lang('challenge_title_thousandth_contribute'), 
                load_lang('challenge_info_thousandth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change("select count(*) from topic where ip = ?"), [ip])
            db_data = curs.fetchall()
            
            disable = 1 if db_data[0][0] >= 1 else 0
            data_html = do_make_challenge_design(
                '💬',
                load_lang('challenge_title_first_discussion'), 
                load_lang('challenge_info_first_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 10 else 0
            data_html = do_make_challenge_design(
                '💡',
                load_lang('challenge_title_tenth_discussion'), 
                load_lang('challenge_info_tenth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 100 else 0
            data_html = do_make_challenge_design(
                '📢',
                load_lang('challenge_title_hundredth_discussion'), 
                load_lang('challenge_info_hundredth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            disable = 1 if db_data[0][0] >= 1000 else 0
            data_html = do_make_challenge_design(
                '📜',
                load_lang('challenge_title_thousandth_discussion'), 
                load_lang('challenge_info_thousandth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
                
            data_html = data_html_green + data_html_red

            disable = 1 if admin_check('all') == 1 else 0
            data_html = do_make_challenge_design(
                '✅',
                load_lang('challenge_title_admin'), 
                load_lang('challenge_info_admin', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
                
            data_html = data_html_green + data_html_red
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('challenge'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = data_html + '''
                    <form method="post">
                        <button id="opennamu_save_button" type="submit">''' + load_lang('reload') + '''</button>
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))