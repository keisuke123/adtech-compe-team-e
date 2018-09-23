from django.db import models
from django.core.cache import cache
import numpy as np
import pandas as pd
from solo.models import SingletonModel
import pickle
import datetime

class CTRPredictor(SingletonModel):
    with open("predictor/models/logisticRegression.pickle", "rb") as f:
        model = pickle.load(f)

    def predict(self, X_new):
        self.model.predict_proba(X_new).T[1]

class ResultsWrapper:

    def __init__(self, shared_dict, ad_id):
        def encode_ad_id(input_dict, ad_id):
            ad_id_map = {"adv02": "advId_2", "adv03": "advId_3", "adv04": "advId_4",
                         "adv05": "advId_5", "adv06": "advId_6", "adv07": "advId_7",
                         "adv08": "advId_8", "adv09": "advId_9","adv10": "advId_10",
                         "adv11": "advId_11", "adv12": "advId_12", "adv13": "advId_13",
                         "adv14": "advId_14", "adv15": "advId_15", "adv16": "advId_16",
                         "adv17": "advId_17", "adv18": "advId_18",
                         "adv19": "advId_19", "adv20": "advId_20"}
            if ad_id in ad_id_map:
                input_dict[ad_id_map[ad_id]] = 1

        self.ad_id = ad_id
        input_dict = shared_dict
        encode_ad_id(input_dict, ad_id)
        df = pd.DataFrame(input_dict, index=[0])
        self.X_new = df.values
        self.ctr = 0

class Preprocessor(SingletonModel):
    with open("predictor/data/media_dict.pickle", "rb") as f:
        media_dict = pickle.load(f)
    with open("predictor/data/click_counts_media_dict.pickle", "rb") as f:
        click_counts_media_dict = pickle.load(f)

    def process(self, feature_dict, ad_id):
        def std_mediaIdCounts(counts):
            mean = 10000.927404
            std = 96.315871
            mediaIdCounts_std = (counts - mean) / std
            return mediaIdCounts_std

        def std_mediaIdClickCounts(Click_counts):
            mean = 2453.497846
            std = 315.193772
            mediaIdClickCounts_std = (Click_counts - mean) / std
            return mediaIdClickCounts_std

        def std_floorPrice(price):
            mean = 8.980794e+03
            std = 2901.756569
            price_std = (price - mean) / std
            return price_std

        def std_age(age):
            mean = 45.500342
            std = 9.377507
            age_std = (age - mean) / std
            return age_std

        def std_income(income):
            mean = 537.125630
            std = 75.535335
            income_std = (income - mean) / std
            return income_std

        def encode_os(ostype):
            if ostype.lower() == 'ios':
                return 1
            else:
                return 0

        def encode_gender(gender, input_dict):
            if gender.lower() == 'male':
                input_dict['male'] = 1
            elif gender.lower() == 'female':
                input_dict['female'] = 1

        def encode_has_child(has_child, input_dict):
            if has_child.lower() == 'yes':
                input_dict['yes'] = 1
            elif has_child.lower() == 'no':
                input_dict['no'] = 1

        def encode_is_married(is_married, input_dict):
            if is_married.lower() == 'yes':
                input_dict['married'] = 1

        def encode_banner_position(banner_position, input_dict):
            pos_dict = {
                2: 'bannerPosition_below',
                3: 'bannerPosition_header',
                4: 'bannerPosition_footer',
                5: 'bannerPosition_Sidebar',
                6: 'bannerPosition_Full'
            }

            if banner_position in pos_dict:
                input_dict[pos_dict[banner_position]] = 1

        def encode_banner_size(banner_size, input_dict):
            size_dict = {
                2: 'bannerSize_2',
                3: 'bannerSize_3',
                4: 'bannerSize_4',
            }
            if banner_size in size_dict:
                input_dictionary[size_dict[banner_size]] = 1

        def encode_timestamp(timestamp, input_dict):
            hour_mapping = {1: 'h_1', 2: 'h_2', 3: 'h_3', 4: 'h_4', 5: 'h_5', 6: 'h_6',
                            7: 'h_7', 8: 'h_8', 9: 'h_9', 10: 'h_10', 11: 'h_11',
                            12: 'h_12', 13: 'h_13', 14: 'h_14', 15: 'h_15', 16: 'h_16',
                            17: 'h_17', 18: 'h_18', 19: 'h_19', 20: 'h_20', 21: 'h_21',
                            22: 'h_22', 23: 'h_23'}

            ts = datetime.datetime.fromtimestamp(timestamp)
            h = ts.hour
            w = ts.weekday()

            if w == 5 | w == 6:
                input_dict['isHoliday'] = 1

            if h in hour_mapping:
                input_dictionary[hour_mapping[h]] = 1

        def encode_media_id(media_id, input_dict):
            input_dict['mediaId_counts'] = std_mediaIdCounts(self.media_dict[media_id])
            input_dict['Click_counts_mediaId'] = std_mediaIdClickCounts(self.click_counts_media_dict[media_id])

        def encode_ad_id(ad_id, input_dict):
            ad_id_map = {"adv02": "advId_2", "adv03": "advId_3", "adv04": "advId_4",
                         "adv05": "advId_5", "adv06": "advId_6", "adv07": "advId_7",
                         "adv08": "advId_8", "adv09": "advId_9", "adv10": "advId_10",
                         "adv11": "advId_11", "adv12": "advId_12", "adv13": "advId_13",
                         "adv14": "advId_14", "adv15": "advId_15", "adv16": "advId_16",
                         "adv17": "advId_17", "adv18": "advId_18",
                         "adv19": "advId_19", "adv20": "advId_20"}
            if ad_id in ad_id_map:
                input_dict[ad_id_map[ad_id]] = 1

        input_dictionary = {
            'floorPrice': 0,
            'isHoliday': 0,
            'bannerPosition_below': 0,
            'bannerPosition_header': 0,
            'bannerPosition_footer': 0,
            'bannerPosition_Sidebar': 0,
            'bannerPosition_Full': 0,
            'os_type_iOS': 0,
            'h_1': 0,
            'h_10': 0,
            'h_11': 0,
            'h_12': 0,
            'h_13': 0,
            'h_14': 0,
            'h_15': 0,
            'h_16': 0,
            'h_17': 0,
            'h_18': 0,
            'h_19': 0,
            'h_2': 0,
            'h_20': 0,
            'h_21': 0,
            'h_22': 0,
            'h_23': 0,
            'h_3': 0,
            'h_4': 0,
            'h_5': 0,
            'h_6': 0,
            'h_7': 0,
            'h_8': 0,
            'h_9': 0,
            'advId_10': 0,
            'advId_11': 0,
            'advId_12': 0,
            'advId_13': 0,
            'advId_14': 0,
            'advId_15': 0,
            'advId_16': 0,
            'advId_17': 0,
            'advId_18': 0,
            'advId_19': 0,
            'advId_2': 0,
            'advId_20': 0,
            'advId_3': 0,
            'advId_4': 0,
            'advId_5': 0,
            'advId_6': 0,
            'advId_7': 0,
            'advId_8': 0,
            'advId_9': 0,
            'mediaId_counts': 0,
            'Click_counts_mediaId': 0,
            'bannerSize_2': 0,
            'bannerSize_3': 0,
            'bannerSize_4': 0,
            'age': 0,
            'income': 0,
            'female': 0,
            'male': 0,
            'not_married': 0,
            'married': 0,
            'no': 0,
            'yes': 0
        }

        input_dictionary['floorPrice'] = std_floorPrice(feature_dict['floorPrice'])
        input_dictionary['age'] = std_age(feature_dict['age'])
        input_dictionary['income'] = std_income(feature_dict['income'])
        input_dictionary['os_type_iOS'] = encode_os(feature_dict['osType'])
        encode_gender(feature_dict['gender'], input_dictionary)
        encode_has_child(feature_dict['hasChild'], input_dictionary)
        encode_is_married(feature_dict['isMarried'], input_dictionary)
        encode_banner_position(feature_dict['bannerPosition'], input_dictionary)
        encode_banner_size(feature_dict['bannerSize'], input_dictionary)
        encode_timestamp(feature_dict['timestamp'], input_dictionary)
        encode_media_id(feature_dict['mediaId'], input_dictionary)
        encode_ad_id(ad_id, input_dictionary)
        df = pd.DataFrame(input_dictionary, index=[0])
        X_new = df.values
        return X_new

