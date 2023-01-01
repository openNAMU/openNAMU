from .tool.func import *

def main_func_easter_egg():
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        random_n = random.randrange(0, 8)
        select_list = [
            'https://www.youtube.com/embed/PWD0ZbR7AOY', # TH06   - Shanghai Teahouse ~ Chinese Tea
            'https://www.youtube.com/embed/HoU29ljOmTE', # TH10.5 - Flawless Clothing of Celestials
            'https://www.youtube.com/embed/PR2vUm-Ald8', # TH06   - U.N. Owen Was Her
            'https://www.youtube.com/embed/opZoEmsu_Lo', # TH09   - Flowering Night
            'https://www.youtube.com/embed/txZFFTusSvw', # TH08   - Reach for the Moon ~ Immortal Smoke
            'https://www.youtube.com/embed/Ixq9xL2tvRU', # TH07   - Phantom Ensemble
            'https://www.youtube.com/embed/-3IAx_r4Au0', # TH17   - Entrusting This World to Idols ~ Idolatrize World
            'https://www.youtube.com/embed/wObZkycA6sc', # TH11   - Last Remote
            # Remix by NyxTheShield
        ]

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change('select name from user_set where id = ? and name = ?'), [ip, 'get_🥚'])
            if not curs.fetchall():
                curs.execute(db_change('insert into user_set (name, id, data) values ("get_🥚", ?, "Y")'), [ip])
                conn.commit()

        return easy_minify(flask.render_template(skin_check(),
            imp = ['Easter Egg', wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '<iframe width="640" height="360" src="' + select_list[random_n] + '" frameborder="0" allowfullscreen></iframe>',
            menu = [['manager', load_lang('return')]]
        ))