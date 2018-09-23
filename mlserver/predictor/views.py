from django.shortcuts import render
import logging
from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework import status
from sklearn.preprocessing import LabelEncoder
import pandas as pd
import numpy as np
from predictor.models import CTRPredictor, Preprocessor, ResultsWrapper
import threading

ad_ids = ["adv01", "adv02", "adv03", "adv04", "adv05", "adv06", "adv07",
          "adv08", "adv09", "adv10", "adv11", "adv12", "adv13", "adv14",
          "adv15", "adv16", "adv17", "adv18", "adv19", "adv20"
          ]


# Create your views here.
@api_view(['POST'])
def ctr(request):
    if request.method == 'POST':
        floor_price = float(),
        media_id = int(),
        timestamp = int(),
        os_type = str(),
        banner_size = int(),
        banner_position = int(),
        gender = str(),
        age = float(),
        income = float(),
        has_child = str()
        is_married = str()
        adv_id = int()

        if 'floorPrice' in request.data:
            floor_price = request.data['floorPrice']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'mediaId' in request.data:
            media_id = request.data['mediaId']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'timestamp' in request.data:
            timestamp = request.data['timestamp']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'osType' in request.data:
            os_type = request.data['osType']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'bannerSize' in request.data:
            banner_size = request.data['bannerSize']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'bannerPosition' in request.data:
            banner_position = request.data['bannerPosition']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'deviceType' in request.data:
            device_type = request.data['deviceType']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'gender' in request.data:
            gender = request.data['gender']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'age' in request.data:
            age = request.data['age']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'income' in request.data:
            income = request.data['income']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'hasChild' in request.data:
            has_child = request.data['hasChild']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'isMarried' in request.data:
            is_married = request.data['isMarried']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)

        processor = Preprocessor.get_solo()
        ctr_predictor = CTRPredictor.get_solo()

        train_dict = {
            'floorPrice': floor_price,  # done
            'mediaId': media_id,  # done
            'timestamp': timestamp,  # done
            'osType': os_type,  # done
            'bannerSize': banner_size,  # done
            'bannerPosition': banner_position,  # done
            'gender': gender,  # done
            'age': age,  # done
            'income': income,  # done
            'hasChild': has_child,  # done
            'isMarried': is_married,  # done
        }

        response_body = {}
        for i in range(20):
            x = processor.process(train_dict, ad_ids[i])
            response_body[ad_ids[i]] = ctr_predictor.model.predict_proba(x).T[1]

        # threads = []
        # for i in range(0, 20, 2):
        #     x1 = processor.process(train_dict, ad_ids[i])
        #     x2 = processor.process(train_dict, ad_ids[i+1])
        #     threads.append(
        #         threading.Thread(target=ctr_predictor.predict, kwargs={'X_new': x1}))
        #     threads.append(
        #         threading.Thread(target=ctr_predictor.predict, kwargs={'X_new': x2}))
        #
        # for t in range(0, 20):
        #     threads[t].start()
        #
        # for t in range(0, 20):
        #     threads[t].join()

        # result_wrappers = []
        # for ad_id in ad_ids:
        #     result_wrappers.append(ResultsWrapper(shared_input_dict, ad_id))

        # for result_wrapper in result_wrappers:
        #     result_wrapper.ctr = ctr_predictor.predict(result_wrapper.X_new)
        #     response_body[result_wrapper.ad_id] = result_wrapper.ctr
        #
        # df = pd.DataFrame(columns=['floorPrice', 'isClick', 'isHoliday',
        #                            'bannerPosition_below', 'bannerPosition_header',
        #                            'bannerPosition_footer', 'bannerPosition_Sidebar',
        #                            'bannerPosition_Full', 'os_type_iOS', 'h_1', 'h_10', 'h_11', 'h_12',
        #                            'h_13', 'h_14', 'h_15', 'h_16', 'h_17', 'h_18', 'h_19', 'h_2', 'h_20',
        #                            'h_21', 'h_22', 'h_23', 'h_3', 'h_4', 'h_5', 'h_6', 'h_7', 'h_8', 'h_9',
        #                            'advId_10', 'advId_11', 'advId_12', 'advId_13', 'advId_14', 'advId_15',
        #                            'advId_16', 'advId_17', 'advId_18', 'advId_19', 'advId_2', 'advId_20',
        #                            'advId_3', 'advId_4', 'advId_5', 'advId_6', 'advId_7', 'advId_8',
        #                            'advId_9', 'mediaId_counts', 'Click_counts_mediaId', 'bannerSize_2',
        #                            'bannerSize_3', 'bannerSize_4', 'age', 'income', 'female', 'male',
        #                            'not_married', 'married', 'no', 'yes'])

        # fit new data into the dataframe

        return Response(response_body, status=status.HTTP_200_OK)
