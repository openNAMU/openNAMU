def start(data):
    data = re.sub('\'\'(?P<in>(?:(?!\'\').)+)\'\'', '<b>\g<in></b>')
    data = re.sub('\'\'\'(?P<in>(?:(?!\'\'\').)+)\'\'\'', '<b>\g<in></b>')

    return data