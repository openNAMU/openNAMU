from .tool.func import *

async def user_count(name = None):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if name == None:
            that = ip_check()
        else:
            that = name

        curs.execute(db_change("select count(*) from history where ip = ?"), [that])
        count = curs.fetchall()
        if count:
            data = count[0][0]
        else:
            data = 0

        curs.execute(db_change("select count(*) from topic where ip = ?"), [that])
        count = curs.fetchall()
        if count:
            data_topic = count[0][0]
        else:
            data_topic = 0
            
        date = get_time()
        date = date.split()
        date = date[0]
        
        data_today = 0
        data_today_len = 0
            
        curs.execute(db_change("select leng from history where date like ? and ip = ?"), [date + '%', that])
        db_data = curs.fetchall()
        for count in db_data:
            count_data = count[0]
            count_data = count_data.replace('+', '')
            count_data = count_data.replace('-', '')

            data_today_len += int(count_data)
            data_today += 1

        date_yesterday = str((
            datetime.datetime.today() + datetime.timedelta(days = -1)
        ).strftime("%Y-%m-%d"))
        
        data_yesterday = 0
        data_yesterday_len = 0
            
        curs.execute(db_change("select leng from history where date like ? and ip = ?"), [date_yesterday + '%', that])
        db_data = curs.fetchall()
        for count in db_data:
            count_data = count[0]
            count_data = count_data.replace('+', '')
            count_data = count_data.replace('-', '')

            data_yesterday_len += int(count_data)
            data_yesterday += 1

        # 한글 지원 필요
        return await render_template(
            await get_lang('count'),
            '''
                <ul>
                    <li><a href="/record/''' + url_pas(that) + '''">''' + await get_lang('edit_record') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/record/topic/''' + url_pas(that) + '''">''' + await get_lang('discussion_record') + '''</a> : ''' + str(data_topic) + '''</a></li>
                    <hr>
                    <li>(''' + await get_lang('beta') + ''') TODAY : ''' + str(data_today) + '''</li>
                    <li>(''' + await get_lang('beta') + ''') TODAY LEN : ''' + str(data_today_len) + '''</li>
                    <li>(''' + await get_lang('beta') + ''') TODAY DIFF : ''' + str(data_today_len - data_yesterday_len) + '''</li>
                    <hr>
                    <li>(''' + await get_lang('beta') + ''') YESTERDAY : ''' + str(data_yesterday) + '''</li>
                    <li>(''' + await get_lang('beta') + ''') YESTERDAY LEN : ''' + str(data_yesterday_len) + '''</li>
                </ul>
            ''',
            0,
            [['user', await get_lang('return')]]
        )
