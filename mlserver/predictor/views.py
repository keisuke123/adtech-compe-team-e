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
        device_type = int(),
        gender = str(),
        age = int(),
        income = int(),
        has_child = str()
        is_married = str()
        id = str()
        adv_id = int()
        device_id = str()

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
        if 'id' in request.data:
            id = request.data['id']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'advId' in request.data:
            id = request.data['id']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)
        if 'deviceId' in request.data:
            device_id = request.data['deviceId']
        else:
            return Response({}, status=status.HTTP_400_BAD_REQUEST)

        train_dict = {
                'deviceId': device_id,
                 'bannerPosition': banner_position,
                 'bannerSize': banner_size,
                 'deviceType': device_type,
                 'floorPrice': floor_price,
                 'id': id,
                 'mediaId': media_id,
                 'osType': os_type,
                 'timestamp': timestamp,
                 'advId': adv_id
            }

        df_train = pd.DataFrame(train_dict, index=[0])

        list_target = list(df_train.drop(["bannerPosition", "bannerSize", "deviceType",
                                          "floorPrice", "mediaId", "timestamp", "advId"], axis=1).columns)
        for target in list_target:
            le = LabelEncoder()
            le.fit(df_train[target])
            df_train[target] = le.transform(df_train[target])
        X_train = df_train.drop(["id"],axis=1).values
        ctr_predictor = CTRPredictor.get_solo()
        #ctr = CTRPredictor.model.predict_proba(X_train[0].reshape(1,-1))

        response_body = [
            {
                'ctr': 0.75
            }
        ]

        return Response(response_body, status=status.HTTP_200_OK)
