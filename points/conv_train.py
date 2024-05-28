import numpy as np
from sklearn.linear_model import LinearRegression
from sklearn.model_selection import train_test_split

# from sklearn.pipeline import make_pipeline
# from sklearn.preprocessing import PolynomialFeatures

# Raw dogging the data like a madman
data = {
    (51.4491681, -3.199995611764706): (2726, 12208),
    (51.46605702659574, -3.273186255319149): (2725, 12213),
    (51.49268334751773, -3.2138919290780144): (2724, 12214),
    (51.49467843615108, -3.181833614208633): (2724, 12215),
    (51.47169133333333, -3.208198): (2725, 12214),
    (51.47624475653595, -3.177102685457516): (2725, 12215),
    (51.45341744736842, -3.265351789473684): (2726, 12213),
    (51.4002662, -3.3372226): (2727, 12212),
    (51.405816151785714, -3.2822252321428573): (2727, 12213),
    (38.54891230136986, -122.80767635616438): (3143, 9493),
    (33.13224931130064, -117.16389652878465): (3296, 9621),
    (33.803229598039216, -116.40195021568627): (3277, 9639),
    (33.77440954340836, -116.46156138263666): (3278, 9637),
    (33.783089, -116.45136071428571): (3278, 9638),
    (33.76915953909844, -116.38718806347747): (3278, 9639),
    (33.74108425, -116.4208725): (3279, 9638),
    (33.7406698483871, -116.39290126774193): (3279, 9639),
    (33.77641966796875, -116.35528003125): (3278, 9640),
    (33.7672240625, -116.2975325625): (3278, 9641),
    (33.746523273654915, -116.35629582374769): (3279, 9640),
    (33.73956386111111, -116.30190438888889): (3279, 9641),
    (33.76454206617647, -116.27252494852941): (3278, 9642),
    (33.73092213636364, -116.23884963636364): (3279, 9642),
    (33.73767665882353, -116.21825970588235): (3279, 9643),
    (33.744890516666665, -116.1865365): (3279, 9644),
}

lat_long = np.array(list(data.keys()))
xy_coords = np.array(list(data.values()))

# Random state is arbitrary, just to make sure we get the same results
lat_long_train, lat_long_test, xy_coords_train, xy_coords_test = train_test_split(
    lat_long, xy_coords, test_size=0.2, random_state=42
)


degrees = 2
model = LinearRegression()

model.fit(lat_long_train, xy_coords_train)
print(model.score(lat_long_test, xy_coords_test))


def lat_long_to_xy(latitude, longitude):
    return model.predict([(latitude, longitude)])[0]


def xy_to_lat_long(x, y):
    return model.inverse_transform([(x, y)])[0]


print(lat_long_to_xy(51.511840, -0.142822))
# Print formula (mx + b)
print(f"y = {model.coef_[0][0]}x + {model.intercept_[0]}")
