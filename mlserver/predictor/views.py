from django.shortcuts import render
import logging
from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework import status
from sklearn.preprocessing import LabelEncoder
import pandas as pd

from predictor.models import CTRPredictor


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

        train_dict = {
            'floorPrice': floor_price, #done
            'mediaId': media_id,
            'timeStamp': timestamp,
            'osType': os_type, #done
            'bannerSize': banner_size, #done
            'bannerPosition': banner_position, #done
            'gender': gender, #done
            'age': age, #done
            'income': income, #done
            'hasChild': has_child, #done
            'isMarried': is_married, #done
            'advId': adv_id
        }

        df = pd.DataFrame(columns=['floorPrice', 'isClick', 'isHoliday',
                                   'bannerPosition_below', 'bannerPosition_header',
                                   'bannerPosition_footer', 'bannerPosition_Sidebar',
                                   'bannerPosition_Full', 'os_type_iOS', 'h_1', 'h_10', 'h_11', 'h_12',
                                   'h_13', 'h_14', 'h_15', 'h_16', 'h_17', 'h_18', 'h_19', 'h_2', 'h_20',
                                   'h_21', 'h_22', 'h_23', 'h_3', 'h_4', 'h_5', 'h_6', 'h_7', 'h_8', 'h_9',
                                   'advId_10', 'advId_11', 'advId_12', 'advId_13', 'advId_14', 'advId_15',
                                   'advId_16', 'advId_17', 'advId_18', 'advId_19', 'advId_2', 'advId_20',
                                   'advId_3', 'advId_4', 'advId_5', 'advId_6', 'advId_7', 'advId_8',
                                   'advId_9', 'mediaId_counts', 'Click_counts_mediaId', 'bannerSize_2',
                                   'bannerSize_3', 'bannerSize_4', 'age', 'income', 'female', 'male',
                                   'not_married', 'married', 'no', 'yes'])

        #fit new data into the dataframe



        #crea

        df_train = pd.DataFrame(train_dict, index=[0])

        list_target = list(df_train.drop(["bannerPosition", "bannerSize", "deviceType",
                                          "floorPrice", "mediaId", "timestamp", "advId"], axis=1).columns)
        for target in list_target:
            le = LabelEncoder()
            le.fit(df_train[target])
            df_train[target] = le.transform(df_train[target])
        X_train = df_train.drop(["id"], axis=1).values
        ctr_predictor = CTRPredictor.get_solo()
        # ctr = CTRPredictor.model.predict_proba(X_train[0].reshape(1,-1))

        response_body = [
            {
                'ctr': 0.75
            }
        ]

        return Response(response_body, status=status.HTTP_200_OK)
