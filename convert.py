# %%
from collections import Counter
import json
import sqlite3

db = sqlite3.connect('testing.db')
cur = db.execute('SELECT id, keywords, awakening_arr FROM monster;')
res = cur.fetchall()

# %%
equivalent: list[tuple[str, int, int, int]] = [
    # annotate_as, small, big, ratio
    ('bound_resist', 10, 52, 2),
    ('skillboost', 21, 56, 2),
    ('sunglasses', 11, 68, 5),
    ('junk_resist', 12, 69, 5),
    ('toxic_resist', 13, 70, 5),
    ('finger', 19, 53, 2),
    ('two_way', 27, 96, 2),
    ('skill_charge', 51, 97, 2),
    ('healing', 9, 98, 2),
    ('fire_orb_plus', 14, 99, 2),
    ('water_orb_plus', 15, 100, 2),
    ('wood_orb_plus', 16, 101, 2),
    ('light_orb_plus', 17, 102, 2),
    ('dark_orb_plus', 18, 103, 2),
    ('heart_orb_plus', 29, 104, 2),
]
output = []
for i, k, a in res:
    aa = Counter(
        v.strip('{}') for v in a.split(' ')
    )
    if not aa['49']:
        continue
    for annotate_as, small, big, ratio in equivalent:
        aa[annotate_as] = aa[str(small)] + aa[str(big)] * ratio
        del aa[str(small)]
        del aa[str(big)]
        if aa[annotate_as] == 0:
            del aa[annotate_as]
    output.append({
        "ID": i,
        "Counter": {
            **aa,
            **Counter(
                'k' + v.strip('{}') for v in k.split(' ')
            ),
        }
    })

with open('data.json', 'w') as f:
    json.dump(output, f)

# %%
