from .tool.func import *

async def bbs_w(bbs_num = '', tool = 'bbs', page = 1, name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        data = ''
        title_name = ''
        sub = ''
        bbs_name_dict = {}

        admin_auth = await acl_check(tool = 'owner_auth')
        admin_auth = 1 if admin_auth == 0 else 0

        if tool == 'record':
            curs.execute(db_change('select set_data, set_id from bbs_set where set_name = "bbs_name"'))
            db_data = curs.fetchall()
            bbs_name_dict = { for_a[1] : for_a[0] for for_a in db_data } if db_data else {}
            
            title_name = name
            sub = '(' + await get_lang('bbs_record') + ')'
            menu = [['user/' + url_pas(name), await get_lang('user_tool')]]
        elif tool == 'comment_record':
            curs.execute(db_change('select set_data, set_id from bbs_set where set_name = "bbs_name"'))
            db_data = curs.fetchall()
            bbs_name_dict = { for_a[1] : for_a[0] for for_a in db_data } if db_data else {}
            
            title_name = name
            sub = '(' + await get_lang('bbs_comment_record') + ')'
            menu = [['user/' + url_pas(name), await get_lang('user_tool')]]
            
        if tool == 'comment_record':
            data += '''
                <table id="main_table_set">
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + await get_lang('editor') + '''</td>
                        <td id="main_table_width">''' + await get_lang('time') + '''</td>
                        <td id="main_table_width">''' + await get_lang('comment') + '''</td>
                    </tr>
            '''
        else:
            data += '''
                <table id="main_table_set">
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + await get_lang('editor') + '''</td>
                        <td id="main_table_width">''' + await get_lang('time') + '''</td>
                        <td id="main_table_width">''' + await get_lang('last_comment_time') + '''</td>
                    </tr>
            '''

        if tool == 'record':
            try:
                curs.execute(db_change('select set_code, set_id, set_data from bbs_data where set_name = "date" and (set_code, set_id) in (select set_code, set_id from bbs_data where set_name = "user_id" and set_data = ?) as sub_query order by set_data desc limit 50'), [name])
            except:
                curs.execute(db_change('select set_code, set_id from bbs_data where set_name = "user_id" and set_data = ? order by set_data desc limit 50'), [name])

            db_data = curs.fetchall()
        elif tool == 'comment_record':
            try:
                curs.execute(db_change('select set_code, set_id, set_data from bbs_data where set_name = "comment_date" and (set_code, set_id) in (select set_code, set_id from bbs_data where set_name = "comment_user_id" and set_data = ?) as sub_query order by set_data desc limit 50'), [name])
            except:
                curs.execute(db_change('select set_code, set_id from bbs_data where set_name = "comment_user_id" and set_data = ? order by set_data desc limit 50'), [name])
            
            db_data = curs.fetchall()

        for for_b in db_data:
            curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where set_code = ? and set_id = ?'), [for_b[0], for_b[1]])
            db_data = curs.fetchall()
            db_data = list(db_data) if db_data else []

            temp_dict = { for_a[0] : for_a[1] for for_a in db_data }

            bbs_name_select = ''
            bbs_split = for_b[1].split('-')
            if tool == 'comment_record':
                bbs_name_select = '(' + bbs_name_dict[bbs_split[0]] + ')'
            elif tool != 'bbs':
                bbs_name_select = '(' + bbs_name_dict[for_b[1]] + ')'

            if tool == 'bbs':
                notice = 1 if len(for_b) > 2 else 0
            else:
                notice = 0

            if tool == 'comment_record':
                curs.execute(db_change('select set_name, set_data, set_code, set_id from bbs_data where set_name = "title" and set_code = ? and set_id = ?'), [bbs_split[1], bbs_split[0]])
                db_data = curs.fetchall()
                db_data = list(db_data) if db_data else []
                for for_a in db_data:
                    temp_dict[for_a[0]] = for_a[1]
            
                comment_link = ''
                if len(bbs_split) > 2:
                    comment_link = '-'.join(bbs_split[2:])
                    
                comment_link += ('-' + for_b[0] if comment_link != '' else for_b[0])
                    
                data += '''
                    <tr>
                        <td>''' + await ip_pas(temp_dict['comment_user_id']) + '''</td>
                        <td>''' + temp_dict['comment_date'] + '''</td>
                        <td>''' + ('#' + comment_link) + '''</td>
                    </tr>
                    <tr>
                        <td colspan="3">
                            <a href="/bbs/w/''' + bbs_split[0] + '/' + bbs_split[1] + '#' + comment_link + '">' + html.escape(temp_dict['title']) + '''</a> 
                            ''' + bbs_name_select + '''
                        </td>
                    </tr>
                '''
            else:
                curs.execute(db_change('select count(*) from bbs_data where set_name = "comment_date" and (set_id = ? or set_id like ?) order by set_code + 0 desc'), [for_b[1] + '-' + for_b[0], for_b[1] + '-' + for_b[0] + '-%'])
                db_data = curs.fetchall()
                comment_count = str(db_data[0][0]) if db_data else '0'

                curs.execute(db_change('select set_data from bbs_data where set_name = "comment_date" and (set_id = ? or set_id like ?) order by set_data desc limit 1'), [for_b[1] + '-' + for_b[0], for_b[1] + '-' + for_b[0] + '-%'])
                db_data = curs.fetchall()
                last_comment_date = db_data[0][0] if db_data else '0'
            
                data += '''
                    <tr class="''' + ('opennamu_comment_color_red' if notice == 1 else '') + '''">
                        <td>''' + await ip_pas(temp_dict['user_id']) + '''</td>
                        <td>''' + temp_dict['date'] + '''</td>
                        <td>''' + last_comment_date + '''</td>
                    </tr>
                    <tr>
                        <td colspan="3">
                            <a href="/bbs/w/''' + for_b[1] + '/' + for_b[0] + '">' + html.escape(temp_dict['title']) + '''</a> 
                            (''' + comment_count + ''') 
                            ''' + bbs_name_select + '''
                        </td>
                    </tr>
                '''
                
        data += '</table>'

        return await render_template(
            title_name,
            data,
            sub,
            menu
        )
