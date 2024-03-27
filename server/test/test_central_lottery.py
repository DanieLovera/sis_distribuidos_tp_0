import os
import unittest
from src.common.central_lottery import *

class TestUtils(unittest.TestCase):
    def setUp(self):
        self.b1 = Bet('0', 'first_0', 'last_0', '10000000','2000-12-20', LOTTERY_WINNER_NUMBER)
        self.b2 = Bet('1', 'first_1', 'last_1', '10000001','2000-12-21', LOTTERY_WINNER_NUMBER + 1)

    def tearDown(self):
        if os.path.exists(STORAGE_FILEPATH):
            os.remove(STORAGE_FILEPATH)

    def test_has_won_with_winner_number_must_be_true(self):
        self.assertTrue(has_won(self.b1))

    def test_has_won_with_loser_number_must_be_false(self):
        self.assertFalse(has_won(self.b2))

    def test_store_bets_and_load_bets_keeps_fields_data(self):
        to_store = [self.b1]
        store_bets(to_store)
        from_load = list(load_bets())

        self.assertEqual(1, len(from_load))
        self._assert_equal_bets(to_store[0], from_load[0])

    def test_store_bets_and_load_bets_keeps_registry_order(self):
        to_store = [self.b1, self.b2]
        store_bets(to_store)
        from_load = list(load_bets())

        self.assertEqual(2, len(from_load))
        self._assert_equal_bets(to_store[0], from_load[0])
        self._assert_equal_bets(to_store[1], from_load[1])

    def _assert_equal_bets(self, b1, b2):
        self.assertEqual(b1.agency, b2.agency)
        self.assertEqual(b1.first_name, b2.first_name)
        self.assertEqual(b1.last_name, b2.last_name)
        self.assertEqual(b1.document, b2.document)
        self.assertEqual(b1.birthdate, b2.birthdate)
        self.assertEqual(b1.number, b2.number)

if __name__ == '__main__':
    unittest.main()
