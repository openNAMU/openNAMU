from .tool.func import *

def bbs_w(bbs_num = '', tool = 'bbs', page = 1, name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        data = ''
        title_name = ''
        sub = ''
        bbs_name_dict = {}

        admin_auth = admin_check()

        if tool == 'bbs':
            curs.execute(db_change('select set_data from bbs_set where set_id = ? and set_name = "bbs_name"'), [bbs_num])
            db_data = curs.fetchall()
            if not db_data:
                return redirect('/bbs/main')
        
            bbs_name = db_data[0][0]
            bbs_num_str = str(bbs_num)

            title_name = bbs_name
            sub = '(' + load_lang('bbs') + ')'
            menu = [['bbs/main', load_lang('return')], ['bbs/edit/' + bbs_num_str, load_lang('add')], ['bbs/set/' + bbs_num_str, load_lang('bbs_set')]]
        elif tool == 'record':
            curs.execute(db_change('select set_data, set_id from bbs_set where set_name = "bbs_name"'))
            db_data = curs.fetchall()
            bbs_name_dict = { for_a[1] : for_a[0] for for_a in db_data } if db_data else {}
            
            title_name = name
            sub = '(' + load_lang('bbs_record') + ')'
            menu = [['user/' + url_pas(name), load_lang('user_tool')]]
        elif tool == 'comment_record':
            curs.execute(db_change('select set_data, set_id from bbs_set where set_name = "bbs_name"'))
            db_data = curs.fetchall()
            bbs_name_dict = { for_a[1] : for_a[0] for for_a in db_data } if db_data else {}
            
            title_name = name
            sub = '(' + load_lang('bbs_comment_record') + ')'
            menu = [['user/' + url_pas(name), load_lang('user_tool')]]
        else:
            curs.execute(db_change('select set_data, set_id from bbs_set where set_name = "bbs_name"'))
            db_data = curs.fetchall()
            if db_data:
                data += '<ul class="opennamu_ul">'
                for for_a in db_data:
                    bbs_name_dict[for_a[1]] = for_a[0]

                    curs.execute(db_change('select set_data from bbs_set where set_name = "bbs_type" and set_id = ?'), [for_a[1]])
                    db_data_2 = curs.fetchall()
                    bbs_type = db_data_2[0][0] if db_data_2 else 'comment'

                    if bbs_type == 'thread':
                        bbs_type = load_lang('thread_base')
                    else:
                        bbs_type = load_lang('comment_base')
                    
                    curs.execute(db_change('select set_data from bbs_data where set_id = ? and set_name = "date" order by set_code + 0 desc limit 1'), [for_a[1]])
                    db_data_2 = curs.fetchall()
                    last_date = ('(' + db_data_2[0][0] + ')') if db_data_2 else ''

                    data += '<li>'
                    data += '<a href="/bbs/w/' + for_a[1] + '">' + html.escape(for_a[0]) + '</a> (' + bbs_type + ') ' + last_date
                    data += '</li>'

                data += '</ul>'
            
            data += '<hr class="main_hr">'

            title_name = load_lang('bbs_main')
            menu = [['other', load_lang('other_tool')]] + ([['bbs/make', load_lang('add')]] if admin_auth == 1 else [])

        if tool == 'comment_record':
            data += '''
                <table id="main_table_set">
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                        <td id="main_table_width">''' + load_lang('comment') + '''</td>
                    </tr>
            '''
        else:
            data += '''
                <table id="main_table_set">
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                        <td id="main_table_width">''' + load_lang('last_comment_time') + '''</td>
                    </tr>
            '''

        if tool == 'bbs':
            curs.execute(db_change('select set_code, set_id, set_name from bbs_data where set_name = "pinned" and set_id like ? order by set_data desc'), [bbs_num])
            db_data = curs.fetchall()
            db_data = list(db_data) if db_data else []
            
            curs.execute(db_change('select set_code, set_id from bbs_data where set_name = "title" and set_id like ? order by set_code + 0 desc'), [bbs_num])
            db_data_2 = curs.fetchall()
            db_data += list(db_data_2) if db_data_2 else []
        elif tool == 'record':
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
        else:
            curs.execute(db_change('select set_code, set_id, set_data from bbs_data where set_name = "date" order by set_data desc limit 50'))
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
                        <td>''' + ip_pas(temp_dict['comment_user_id']) + '''</td>
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
                        <td>''' + ip_pas(temp_dict['user_id']) + '''</td>
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

        return easy_minify(flask.render_template(skin_check(),
            imp = [title_name, wiki_set(), wiki_custom(), wiki_css([sub, 0])],
            data = data,
            menu = menu
        ))