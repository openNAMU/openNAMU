from .tool.func import *

async def bbs_w_comment_tool(bbs_num = '', post_num = '', comment_num = ''):
    with get_db_connect() as conn:
        data = ''
        
        bbs_num_str = str(bbs_num)
        post_num_str = str(post_num)
        
        data += '''
            <h2>''' + await get_lang('tool') + '''</h2>
            <ul>
                <li><a href="/bbs/raw/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + await get_lang('raw') + '''</a></li>
                <li><a href="/bbs/edit/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + await get_lang('edit') + '''</a></li>
            </ul>
        '''

        if await acl_check('', 'owner_auth', '', '') != 1:
            data += '''
                <h3>''' + await get_lang('owner') + '''</h2>
                <ul>
                    <li><a href="/bbs/delete/''' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '/' + url_pas(comment_num) + '">' + await get_lang('delete') + '''</a></li>
                </ul>
            '''

        return await render_template(
            await get_lang('bbs_comment_tool'),
            data,
            0,
            [['bbs/w/' + url_pas(bbs_num_str) + '/' + url_pas(post_num_str) + '#' + url_pas(comment_num), await get_lang('return')]]
        )
