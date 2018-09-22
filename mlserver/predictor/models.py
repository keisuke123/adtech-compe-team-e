from django.db import models
from django.core.cache import cache
import numpy as np
import pandas as pd
from solo.models import SingletonModel
import pickle

class CTRPredictor(SingletonModel):
	model = 'a'
    #with open("predictor/models/test_model.pickle", "rb") as f:
     #   model = pickle.load(f)
