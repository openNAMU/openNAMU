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

async def user_challenge():    
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        ip = ip_check()
        if ip_or_user(ip) == 1:
            return redirect(conn, '/user')

        if flask.request.method == 'POST':
            user_exp = 0

            curs.execute(db_change('select count(*) from history where ip = ?'), [ip])
            db_data = curs.fetchall()
            if not db_data:
                db_data = [[0]]

            user_exp += 5 * db_data[0][0]

            if db_data[0][0] >= 1:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_first_contribute'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_first_contribute', ?, '1')"), [ip])
                user_exp += 500

            if db_data[0][0] >= 10:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_tenth_contribute'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_tenth_contribute', ?, '1')"), [ip])
                user_exp += 1000

            if db_data[0][0] >= 100:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_hundredth_contribute'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_hundredth_contribute', ?, '1')"), [ip])
                user_exp += 3000        

            if db_data[0][0] >= 1000:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_thousandth_contribute'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_thousandth_contribute', ?, '1')"), [ip])
                user_exp += 10000

            curs.execute(db_change("select count(*) from topic where ip = ?"), [ip])
            db_data = curs.fetchall()
            if not db_data:
                db_data = [[0]]

            user_exp += 5 * db_data[0][0]

            if db_data[0][0] >= 1:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_first_discussion'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_first_discussion', ?, '1')"), [ip])
                user_exp += 500    

            if db_data[0][0] >= 10:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_tenth_discussion'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_tenth_discussion', ?, '1')"), [ip])
                user_exp += 1000

            if db_data[0][0] >= 100:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_hundredth_discussion'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_hundredth_discussion', ?, '1')"), [ip])
                user_exp += 3000

            if db_data[0][0] >= 1000:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_thousandth_discussion'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_thousandth_discussion', ?, '1')"), [ip])
                user_exp += 10000        

            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_admin', ip])
            db_data = curs.fetchall()
            if await acl_check(tool = 'all_admin_auth') != 1 or db_data:
                curs.execute(db_change("delete from user_set where id = ? and name = 'challenge_admin'"), [ip])
                curs.execute(db_change("insert into user_set (name, id, data) values ('challenge_admin', ?, '1')"), [ip])
                user_exp += 10000

            exp = user_exp
            level = 0
            while 1:
                if exp >= (500 + level * 50):
                    exp -= (500 + level * 50)
                    level += 1
                else:
                    break

            curs.execute(db_change("delete from user_set where id = ? and name = 'level'"), [ip])
            curs.execute(db_change("insert into user_set (name, id, data) values ('level', ?, ?)"), [ip, level])

            curs.execute(db_change("delete from user_set where id = ? and name = 'experience'"), [ip])
            curs.execute(db_change("insert into user_set (name, id, data) values ('experience', ?, ?)"), [ip, exp])

            return redirect(conn, '/challenge')
        else:
            data_html_green = ''
            data_html_red = ''
            
            data_html_green += do_make_challenge_design(
                '🌳',
                get_lang(conn, 'challenge_title_register'), 
                get_lang(conn, 'challenge_info_register', 1),
                1
            )
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_first_contribute', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '🔰',
                get_lang(conn, 'challenge_title_first_contribute'), 
                get_lang(conn, 'challenge_info_first_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_tenth_contribute', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '📝',
                get_lang(conn, 'challenge_title_tenth_contribute'), 
                get_lang(conn, 'challenge_info_tenth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_hundredth_contribute', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '🖊️',
                get_lang(conn, 'challenge_title_hundredth_contribute'), 
                get_lang(conn, 'challenge_info_hundredth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_thousandth_contribute', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '🏅',
                get_lang(conn, 'challenge_title_thousandth_contribute'), 
                get_lang(conn, 'challenge_info_thousandth_contribute', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_first_discussion', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '💬',
                get_lang(conn, 'challenge_title_first_discussion'), 
                get_lang(conn, 'challenge_info_first_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_tenth_discussion', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '💡',
                get_lang(conn, 'challenge_title_tenth_discussion'), 
                get_lang(conn, 'challenge_info_tenth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_hundredth_discussion', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '📢',
                get_lang(conn, 'challenge_title_hundredth_discussion'), 
                get_lang(conn, 'challenge_info_hundredth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
            
            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_thousandth_discussion', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '📜',
                get_lang(conn, 'challenge_title_thousandth_discussion'), 
                get_lang(conn, 'challenge_info_thousandth_discussion', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
                
            data_html = data_html_green + data_html_red

            curs.execute(db_change('select data from user_set where name = ? and id = ?'), ['challenge_admin', ip])
            db_data = curs.fetchall()
            disable = 1 if db_data else 0
            data_html = do_make_challenge_design(
                '☑️',
                get_lang(conn, 'challenge_title_admin'), 
                get_lang(conn, 'challenge_info_admin', 1),
                disable
            )
            if disable == 1:
                data_html_green += data_html
            else:
                data_html_red += data_html
                
            data_html = data_html_green + data_html_red
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'challenge_and_level_manage'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = data_html + '''
                    <form method="post">
                        <div id="opennamu_get_user_info">''' + html.escape(ip) + '''</div>
                        <hr class="main_hr">
                        <button id="opennamu_save_button" type="submit">''' + get_lang(conn, 'reload') + '''</button>
                    </form>
                ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))