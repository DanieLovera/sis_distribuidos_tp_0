import unittest
from src.common.bet import *

class TestUtils(unittest.TestCase):
    def test_bet_init_must_keep_fields(self):
        b = Bet('1', 'first', 'last', '10000000','2000-12-20', 7500)
        self.assertEqual(1, b.agency)
        self.assertEqual('first', b.first_name)
        self.assertEqual('last', b.last_name)
        self.assertEqual('10000000', b.document)
        self.assertEqual(datetime.date(2000, 12, 20), b.birthdate)
        self.assertEqual(7500, b.number)

if __name__ == '__main__':
    unittest.main()
