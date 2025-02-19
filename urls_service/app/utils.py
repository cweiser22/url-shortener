import secrets
import string


def generate_short_code() -> str:
    short_url = ''.join(secrets.choice(string.ascii_letters + string.digits) for _ in range(7))
    return short_url


def serialize_bson_document(document: dict):
    if document and '_id' in document:
        document['_id'] = str(document['_id'])
    return document
