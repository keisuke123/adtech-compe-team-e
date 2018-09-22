from django.db import models
from django.core.cache import cache
import numpy as np
import pandas as pd
from solo.models import SingletonModel
import pickle


class CTRPredictor(SingletonModel):
    with open("predictor/models/test_model.pickle", "rb") as f:
        model = pickle.load(f)


class Preprocessor(SingletonModel):

    def process(feature_dict):

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
            def switch(pos):
                return pos_dict[pos]

            if banner_position in pos_dict:
                input_dict[switch(banner_position)] = 1

        def encode_banner_size(banner_size, input_dict):
            size_dict = {
                2: 'bannerSize_2',
                3: 'bannerSize_3',
                4: 'bannerSize_4',
            }
            def switch(size):
                return size_dict[size]

            if banner_size in size_dict:
                input_dictionary[switch(banner_size)] = 1



        input_dictionary = {
            'floorPrice': 0,
            'isClick': 0,
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
        input_dictionary['os_type_iOS'] = encode_os(feature_dict['os_type'])
        encode_gender(feature_dict['gender'], input_dictionary)
        encode_has_child(feature_dict['hasChild'], input_dictionary)
        encode_is_married(feature_dict['isMarried', input_dictionary])
        encode_banner_position(feature_dict['bannerPosition'], input_dictionary)
        encode_banner_size(feature_dict['bannerSize'], input_dictionary)

        return input_dictionary
